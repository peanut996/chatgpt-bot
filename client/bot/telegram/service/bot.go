package service

import (
	"chatgpt-bot/bot/telegram/handler"
	"chatgpt-bot/bot/telegram/limiter"
	"chatgpt-bot/cfg"
	"chatgpt-bot/constant/cmd"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/db"
	"chatgpt-bot/engine"
	"chatgpt-bot/model"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Bot struct {
	config *cfg.Config
	tgBot  *tgbotapi.BotAPI
	engine engine.Engine

	chatTaskChannel chan model.ChatTask
	gpt4TaskChannel chan model.ChatTask
	handlers        map[handler.BotCmd]handler.CommandHandler
	limiters        []limiter.Limiter
	gpt4Limiters    []limiter.Limiter
}

func (b *Bot) SelfID() int64 {
	return b.tgBot.Self.ID
}

func (b *Bot) Config() *cfg.Config {
	return b.config
}

func (b *Bot) TGBot() *tgbotapi.BotAPI {
	return b.tgBot
}

func (b *Bot) Init(cfg *cfg.Config) error {
	b.config = cfg
	if cfg.GPT4Limiter == nil || cfg.GPT3Limiter == nil {
		return errors.New(botError.MissingRequiredConfig + ": limiters")
	}

	if utils.IsAnyStringEmpty(
		b.config.TelegramChannelName, b.config.TelegramGroupName) {
		return errors.New(botError.MissingRequiredConfig + ": group and channel name")
	}

	d := db.NewSQLiteDB()
	err := d.Init(cfg)
	if err != nil {
		return err
	}

	userRepository := repository.NewUserRepository(d)
	userInviteRecordRepository := repository.NewUserInviteRecordRepository(d)

	bot, err := tgbotapi.NewBotAPI(cfg.BotConfig.TelegramBotToken)
	if err != nil {
		return err
	}
	b.tgBot = bot
	b.engine = engine.GetEngine(cfg.EngineConfig.EngineType)
	err = b.engine.Init(cfg)
	if err != nil {
		return err
	}

	b.chatTaskChannel = make(chan model.ChatTask, 100)
	b.gpt4TaskChannel = make(chan model.ChatTask, 100)

	b.handlers = make(map[handler.BotCmd]handler.CommandHandler)

	b.registerCommandHandler(
		handler.NewStartCommandHandler(userRepository, userInviteRecordRepository),
		handler.NewPingCommandHandler(), handler.NewPprofCommandHandler(), handler.NewLimiterCommandHandler(),
		handler.NewInviteCommandHandler(userRepository),
		handler.NewCountCommandHandler(userRepository),
		handler.NewQueryCommandHandler(userRepository, userInviteRecordRepository),
		handler.NewDonateCommandHandler(),
		handler.NewStatusCommandHandler(userRepository, userInviteRecordRepository),
		handler.NewPushCommandHandler(userRepository),
		handler.NewVIPCommandHandler(userRepository),
		handler.NewAccessCommandHandler(userRepository, userInviteRecordRepository, cfg.BotConfig.SALT),
	)
	initLimiters(cfg, b, userRepository, userInviteRecordRepository)

	go b.loopAndFinishChatTask()

	log.Printf("[Init] telegram bot init success, bot name: %s", b.tgBot.Self.UserName)
	return nil
}

func initLimiters(cfg *cfg.Config, b *Bot, userRepository *repository.UserRepository, recordRepository *repository.UserInviteRecordRepository) {
	common := limiter.NewCommonMessageLimiter()
	singleton := limiter.NewSingletonMessageLimiter()
	join := limiter.NewJoinMessageLimiter()
	invite := limiter.NewInviteCountLimiter(userRepository, recordRepository)
	count := limiter.NewRemainCountMessageLimiter(userRepository)
	user := limiter.NewUserLimiter(userRepository)
	b.registerGPT3Limiter(common, singleton, user)
	b.registerGPT4Limiter(common, singleton, user)

	if cfg.GPT3Limiter.Join {
		b.registerGPT3Limiter(join)
	}
	if cfg.GPT3Limiter.RemainCount {
		b.registerGPT3Limiter(count)
	}

	b.registerGPT3Limiter(
		limiter.NewRateLimiter(cfg.GPT3Limiter.Capacity, cfg.GPT3Limiter.Duration, false,
			userRepository, recordRepository))

	if cfg.GPT4Limiter.Join {
		b.registerGPT4Limiter(join)
	}

	if cfg.GPT4Limiter.Invite {
		b.registerGPT4Limiter(invite)
	}

	if cfg.GPT3Limiter.RemainCount {
		b.registerGPT4Limiter(count)
	}

	b.registerGPT4Limiter(
		limiter.NewRateLimiter(cfg.GPT4Limiter.Capacity, cfg.GPT4Limiter.Duration, true,
			userRepository, recordRepository))

}

func NewTelegramBot() *Bot {
	return &Bot{}
}

func (b *Bot) Run() {
	log.Println("[Run] start telegram bot")
	go b.fetchUpdates()
}

func (b *Bot) fetchUpdates() {
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60
	config.AllowedUpdates = []string{"message", "edited_message", "channel_post", "edited_channel_post", "chat_member"}

	botChannel := b.tgBot.GetUpdatesChan(config)
	for {
		select {
		case update, ok := <-botChannel:
			if !ok {
				b.tgBot.StopReceivingUpdates()
				botChannel = b.tgBot.GetUpdatesChan(config)
				log.Println("[FetchUpdates] channel closed, fetch again")
				continue
			}
			go b.handleUpdate(update)
		case <-time.After(30 * time.Second):
		}
	}
}

func (b *Bot) loopAndFinishChatTask() {
	for {
		select {
		case task := <-b.gpt4TaskChannel:
			go b.finishChatTask(task)
		case task := <-b.chatTaskChannel:
			b.finishChatTask(task)
		case <-time.After(30 * time.Second):
		}

	}
}

func (b *Bot) finishChatTask(task model.ChatTask) {
	log.Printf("[finishChatTask] start chat task %s", task.String())
	b.logToChannel(task.GetFormattedQuestion())
	b.sendTyping(task.Chat)

	chatCtx := model.NewChatContext(task.Question, utils.Int64ToString(task.From), "")
	if task.IsGPT4Message {
		chatCtx.Model = "gpt-4"
	}
	res, err := b.engine.Chat(chatCtx)
	if err != nil {
		task.Answer = err.Error()
	} else {
		task.Answer = res
	}
	b.sendTyping(task.Chat)
	b.sendFromChatTask(task)
	b.logToChannel(task.GetFormattedAnswer())

	b.runLimitersCallBack(task.GetRawMessage(), true)

	log.Printf("[finishChatTask] end chat task: %s", task.String())
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[Update] 【msg】:%s", utils.ToJson(update.Message))

	log.Printf("[Update] 【type】:%s, 【from】:%s 【text】: %s",
		update.Message.Chat.Type, model.From(update.Message.From).String(),
		update.Message.Text)

	if update.Message.IsCommand() && !IsGPTMessage(*update.Message) {
		b.execCommand(*update.Message)
		return
	}

	if IsGPTMessage(*update.Message) && strings.Trim(update.Message.CommandArguments(), " ") == "" {
		b.SafeReplyMsg(update.Message.Chat.ID, update.Message.MessageID, fmt.Sprintf(tip.GPTLackTextTipTemplate,
			update.Message.Command(), update.Message.Command()))
		return
	}

	b.handleMessage(*update.Message)

}

func (b *Bot) handleMessage(message tgbotapi.Message) {

	ok := b.checkLimiters(message)
	if !ok {
		b.runLimitersCallBack(message, false)
		return
	}

	if !IsGPT4Message(message) {
		b.sendQueueToast(message.Chat.ID, message.MessageID)
	}

	b.publishChatTask(message)

}

func (b *Bot) publishChatTask(message tgbotapi.Message) {
	log.Printf("[publishChatTask] with message %s", utils.ToJson(message))
	chatTask := model.NewChatTask(message)
	user, err := b.GetUserInfo(message.From.ID)
	if err == nil {
		chatTask.User = user
	}
	if IsGPT4Message(message) && !b.config.Downgrade {
		b.gpt4TaskChannel <- *chatTask
	} else {
		chatTask.IsGPT4Message = false
		b.chatTaskChannel <- *chatTask
	}
	b.sendTyping(chatTask.Chat)
}

func (b *Bot) registerCommandHandler(handlers ...handler.CommandHandler) {
	for _, commandHandler := range handlers {
		b.handlers[commandHandler.Cmd()] = commandHandler
	}
}

func (b *Bot) execCommand(message tgbotapi.Message) {
	command := message.Command()
	if !cmd.IsBotCmd(command) {
		return
	}
	commandHandler, ok := b.handlers[command]
	if !ok {
		b.SafeSend(tgbotapi.NewMessage(message.Chat.ID, tip.UnknownCmdTip))
		return
	}

	err := commandHandler.Run(b, message)
	if err != nil {
		log.Println("[CommandHandler]exec handler encounter error: " + err.Error())
		b.SafeReplyMsg(message.Chat.ID, message.MessageID, botError.InternalError)
	}
}

func (b *Bot) registerLimiter(isGPT4 bool, limiters ...limiter.Limiter) {
	if isGPT4 {
		b.gpt4Limiters = append(b.gpt4Limiters, limiters...)
		return
	}
	b.limiters = append(b.limiters, limiters...)

}

func (b *Bot) registerGPT3Limiter(limiters ...limiter.Limiter) {
	b.registerLimiter(false, limiters...)
}
func (b *Bot) registerGPT4Limiter(limiters ...limiter.Limiter) {
	b.registerLimiter(true, limiters...)
}

func (b *Bot) checkLimiters(m tgbotapi.Message) bool {
	limiters := b.limiters
	if IsGPTMessage(m) && m.Command() == cmd.GPT4 {
		limiters = b.gpt4Limiters
	}
	for _, l := range limiters {
		ok, err := l.Allow(b, m)
		if !ok {
			if utils.IsNotEmpty(err) {
				log.Printf("[CheckLimiter] limiter encounter type: %s error: %s", reflect.TypeOf(l).String(), err)
				b.SafeReplyMsg(m.Chat.ID, m.MessageID, err)
			}
			return false
		}
	}
	return true
}

func (b *Bot) runLimitersCallBack(m tgbotapi.Message, success bool) {
	limiters := b.limiters
	if IsGPTMessage(m) && m.Command() == cmd.GPT4 {
		limiters = b.gpt4Limiters
	}
	for _, l := range limiters {
		l.CallBack(b, m, success)
	}
}
