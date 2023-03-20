package telegram

import (
	"chatgpt-bot/cfg"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/db"
	"chatgpt-bot/engine"
	"chatgpt-bot/model"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Bot struct {
	tgBot  *tgbotapi.BotAPI
	engine engine.Engine

	chatTaskChannel chan model.ChatTask
	maxQueueChannel chan interface{}

	groupName     string
	channelName   string
	limitGroup    bool
	limitPrivate  bool
	logChannelID  int64
	enableLimiter bool
	admin         int64

	handlers     map[BotCmd]CommandHandler
	limiters     []Limiter
	gpt4Limiters []Limiter
}

func (b *Bot) Init(cfg *cfg.Config) error {
	b.channelName = cfg.BotConfig.TelegramChannelName
	b.groupName = cfg.BotConfig.TelegramGroupName
	b.limitPrivate = cfg.BotConfig.ShouldLimitPrivate
	b.limitGroup = cfg.BotConfig.ShouldLimitGroup
	b.logChannelID = cfg.BotConfig.LogChannelID
	b.admin = cfg.BotConfig.AdminID

	if utils.IsAnyStringEmpty(
		b.channelName, b.groupName) {
		return errors.New(botError.MissingRequiredConfig)
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

	b.chatTaskChannel = make(chan model.ChatTask, 1)
	b.maxQueueChannel = make(chan interface{}, 3)
	b.handlers = make(map[BotCmd]CommandHandler)

	b.enableLimiter = cfg.BotConfig.RateLimiterConfig.Enable

	b.registerCommandHandler(
		NewStartCommandHandler(userRepository, userInviteRecordRepository),
		NewPingCommandHandler(), NewPprofCommandHandler(), NewLimiterCommandHandler(),
		NewInviteCommandHandler(userRepository),
		NewCountCommandHandler(userRepository),
		NewChatCommandHandler(),
		NewQueryCommandHandler(userRepository, userInviteRecordRepository),
		NewDonateCommandHandler(),
	)
	initLimiters(cfg, b, userRepository)

	go b.loopAndFinishChatTask()

	log.Printf("[Init] telegram bot init success, bot name: %s", b.tgBot.Self.UserName)
	return nil
}

func initLimiters(_ *cfg.Config, b *Bot, userRepository *repository.UserRepository) {
	common := NewCommonMessageLimiter()
	singleton := NewSingletonMessageLimiter()
	join := NewJoinMessageLimiter()
	user := NewUserLimiter(userRepository)

	b.registerGPT3Limiter(common, singleton, user,
		NewRateLimiter(1, 60),
	)

	b.registerGPT4Limiter(
		common, singleton, join, user,
		NewRemainCountMessageLimiter(userRepository),
		NewRateLimiter(1, 300),
	)
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
		case task := <-b.chatTaskChannel:
			log.Println("[LoopAndFinishChatTask] got a task to finishChatTask")
			go b.finishChatTask(task)
		case <-time.After(30 * time.Second):
		}

	}
}

func (b *Bot) finishChatTask(task model.ChatTask) {
	b.maxQueueChannel <- struct{}{}
	defer func() {
		<-b.maxQueueChannel
	}()

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
	log.Printf("[Update] 【chat】:%s, 【from】:%s, 【msg】:%s", utils.ToJson(update.Message.Chat),
		utils.ToJson(update.Message.From),
		utils.ToJson(update.Message))

	if update.Message.IsCommand() && !IsGPT4Message(*update.Message) {
		b.execCommand(*update.Message)
		return
	}

	if IsGPT4Message(*update.Message) && strings.Trim(update.Message.CommandArguments(), " ") == "" {
		b.safeSendMsg(update.Message.Chat.ID, tip.GPT4LackTextTip)
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
		b.publishChatTask(message, false)
		return
	}
	b.publishChatTask(message, true)

}

func (b *Bot) publishChatTask(message tgbotapi.Message, isGPT4Task bool) {
	log.Printf("[publishChatTask] with message %s", utils.ToJson(message))
	chatTask := &model.ChatTask{}
	if isGPT4Task {
		chatTask = model.NewGPT4ChatTask(message)
	} else {
		chatTask = model.NewChatTask(message)
	}
	user, err := b.getUserInfo(message.From.ID)
	if err == nil {
		chatTask.User = user
	}
	b.chatTaskChannel <- *chatTask
	b.sendTyping(chatTask.Chat)
}

func (b *Bot) registerCommandHandler(handlers ...CommandHandler) {
	for _, handler := range handlers {
		b.handlers[handler.Cmd()] = handler
	}
}

func (b *Bot) execCommand(message tgbotapi.Message) {
	cmd := message.Command()
	handler, ok := b.handlers[cmd]
	if !ok {
		b.safeSend(tgbotapi.NewMessage(message.Chat.ID, tip.UnknownCmdTip))
		return
	}

	err := handler.Run(b, message)
	if err != nil {
		log.Println("exec handler encounter error: " + err.Error())
		b.safeReplyMsg(message.Chat.ID, message.MessageID, botError.InternalError)
	}
}

func (b *Bot) registerLimiter(isGPT4 bool, limiters ...Limiter) {
	if isGPT4 {
		b.gpt4Limiters = append(b.gpt4Limiters, limiters...)
		return
	}
	b.limiters = append(b.limiters, limiters...)

}

func (b *Bot) registerGPT3Limiter(limiters ...Limiter) {
	b.registerLimiter(false, limiters...)
}
func (b *Bot) registerGPT4Limiter(limiters ...Limiter) {
	b.registerLimiter(true, limiters...)
}

func (b *Bot) checkLimiters(m tgbotapi.Message) bool {
	limiters := b.limiters
	if IsGPT4Message(m) {
		limiters = b.gpt4Limiters
	}
	for _, limiter := range limiters {
		ok, err := limiter.Allow(b, m)
		if !ok {
			if utils.IsNotEmpty(err) {
				log.Printf("[CheckLimiter] limiter encounter type: %s error: %s", reflect.TypeOf(limiter).String(), err)
				b.sendErrorMessage(m.Chat.ID, m.MessageID, err)
			}
			return false
		}
	}
	return true
}

func (b *Bot) runLimitersCallBack(m tgbotapi.Message, success bool) {
	limiters := b.limiters
	if IsGPT4Message(m) {
		limiters = b.gpt4Limiters
	}
	for _, limiter := range limiters {
		limiter.CallBack(b, m, success)
	}
}
