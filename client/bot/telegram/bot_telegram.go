package telegram

import (
	"chatgpt-bot/cfg"
	"chatgpt-bot/constant"
	"chatgpt-bot/db"
	"chatgpt-bot/engine"
	"chatgpt-bot/model"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"errors"
	"log"
	"strconv"
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

	handlers map[BotCmd]CommandHandler
	limiters []MessageLimiter
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
		return errors.New(constant.MissingRequiredConfig)
	}

	d := db.NewSQLiteDB()
	err := d.Init(cfg)
	if err != nil {
		return err
	}

	userRepository := repository.NewUserRepository(d)

	bot, err := tgbotapi.NewBotAPI(cfg.BotConfig.TelegramBotToken)
	if err != nil {
		return err
	}
	b.tgBot = bot
	b.tgBot.Debug = true
	b.engine = engine.GetEngine(cfg.EngineConfig.EngineType)
	err = b.engine.Init(cfg)
	if err != nil {
		return err
	}

	b.chatTaskChannel = make(chan model.ChatTask, 1)
	b.maxQueueChannel = make(chan interface{}, 3)

	b.enableLimiter = cfg.BotConfig.RateLimiterConfig.Enable

	b.registerCommandHandler(NewStartCommand(), NewPingCommand(), NewPprofCommand(), NewLimiterCommand())
	b.registerLimiter(NewCommonMessageLimiter(),
		NewSingleMessageLimiter(),
		NewPrivateMessageLimiter(userRepository),
		NewRateLimiter(cfg.BotConfig.RateLimiterConfig.Capacity, cfg.BotConfig.RateLimiterConfig.Duration),
	)

	go b.loopAndFinishChatTask()

	log.Printf("[Init] telegram bot init success, bot name: %s", b.tgBot.Self.UserName)
	return nil
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
	res, err := b.engine.Chat(task.Question, strconv.FormatInt(task.From, 10))
	if err != nil {
		task.Answer = err.Error()
	} else {
		task.Answer = res
	}
	b.sendTyping(task.Chat)
	b.sendFromChatTask(task)
	b.logToChannel(task.GetFormattedAnswer())
	b.runLimitersCallBack(task.GetRawMessage())
	log.Printf("[finishChatTask] end chat task: %s", task.String())
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[Update] 【chat】:%s, 【from】:%s, 【msg】:%s", utils.ToJson(update.Message.Chat),
		utils.ToJson(update.Message.From),
		utils.ToJson(update.Message))
	if update.Message.IsCommand() {
		b.execCommand(update.Message.Command(), update)
	}

	if update.Message != nil {
		b.handleMessage(update.Message)
	}

}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	ok := b.checkLimiters(*message)
	if !ok {
		return
	}
	b.publishChatTask(*message)
}

func (b *Bot) publishChatTask(message tgbotapi.Message) {
	log.Printf("[publishChatTask] with message %s", utils.ToJson(message))
	chatTask := model.NewChatTask(message)
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

func (b *Bot) execCommand(cmd string, update tgbotapi.Update) {
	handler, ok := b.handlers[cmd]
	if !ok {
		b.safeSend(tgbotapi.NewMessage(update.Message.Chat.ID, constant.UnknownCmdTip))
	}

	err := handler.Run(b, update)
	if err != nil {
		log.Println("exec handler encounter error: " + err.Error())
		b.safeSend(tgbotapi.NewMessage(update.Message.Chat.ID, constant.InternalError))
	}
}

func (b *Bot) registerLimiter(limiters ...MessageLimiter) {
	b.limiters = append(b.limiters, limiters...)
}

func (b *Bot) checkLimiters(m tgbotapi.Message) bool {
	for _, limiter := range b.limiters {
		ok, err := limiter.Allow(b, m)
		if !ok {
			if utils.IsNotEmpty(err) {
				log.Println("[CheckLimiter] limiter encounter error: " + err)
				b.sendErrorMessage(m.Chat.ID, m.MessageID, err)
			}
			return false
		}
	}
	return true
}

func (b *Bot) runLimitersCallBack(m tgbotapi.Message) {
	for _, limiter := range b.limiters {
		limiter.CallBack(b, m)
	}
}
