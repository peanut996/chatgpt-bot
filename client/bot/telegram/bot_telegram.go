package telegram

import (
	"chatgpt-bot/cfg"
	"chatgpt-bot/constant"
	"chatgpt-bot/db"
	"chatgpt-bot/engine"
	"chatgpt-bot/middleware"
	"chatgpt-bot/model"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Bot struct {
	tgBot          *tgbotapi.BotAPI
	engine         engine.Engine
	session        *sync.Map
	limiter        *middleware.Limiter
	db             db.BotDB
	userRepository *repository.UserRepository

	taskChan     chan *model.ChatTask
	maxQueueChan chan interface{}

	groupName     string
	channelName   string
	limitGroup    bool
	limitPrivate  bool
	logChat       int64
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
	b.logChat = cfg.BotConfig.LogChannelID
	b.admin = cfg.BotConfig.AdminID

	if utils.IsAnyStringEmpty(
		b.channelName, b.groupName) {
		return errors.New(constant.MissingRequiredConfig)
	}

	db := db.NewSQLiteDB()
	err := db.Init(cfg)
	if err != nil {
		return err
	}
	b.db = db

	b.userRepository = repository.NewUserRepository(b.db)

	b.session = &sync.Map{}
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

	b.taskChan = make(chan *model.ChatTask, 1)
	b.maxQueueChan = make(chan interface{}, 3)

	b.enableLimiter = cfg.BotConfig.RateLimiterConfig.Enable
	b.limiter = middleware.NewLimiter(cfg.BotConfig.RateLimiterConfig.Capacity,
		cfg.BotConfig.RateLimiterConfig.Duration)

	b.registerCommandHandler(NewStartCommand(), NewPingCommand(), NewPprofCommand(), NewLimiterCommand())
	b.registerLimiter(NewCommonMessageLimiter())

	go b.loopAndFinishChatTask()

	log.Printf("[Init] telegram bot init success, bot name: %s", b.tgBot.Self.UserName)
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
		case task := <-b.taskChan:
			log.Println("[LoopAndFinishChatTask] got a task to finish")
			go b.Finish(task)
		case <-time.After(30 * time.Second):
		}

	}
}

func (b *Bot) Finish(task *model.ChatTask) {
	b.maxQueueChan <- struct{}{}
	defer func() {
		<-b.maxQueueChan
	}()
	log.Printf("[Finish] start chat task %s", task.String())
	defer b.session.Delete(task.From)
	b.Log(task.GetFormattedQuestion())
	b.sendTyping(task)
	res, err := b.engine.Chat(task.Question, strconv.FormatInt(task.From, 10))
	if err != nil {
		task.Answer = err.Error()
	} else {
		task.Answer = res
	}
	b.sendTyping(task)
	b.Send(task)
	b.Log(task.GetFormattedAnswer())

	go b.userRepository.DecreaseCount(strconv.FormatInt(task.From, 10))
	log.Printf("[Finish] end chat task: %s", task.String())
}

func (b *Bot) Log(log string) {
	go func(s string) {
		msg := tgbotapi.NewMessage(b.logChat, s)
		msg.ParseMode = tgbotapi.ModeMarkdown
		b.safeSend(msg)
	}(log)
}

func (b *Bot) Send(task *model.ChatTask) {
	msg := tgbotapi.NewMessage(task.Chat, task.Question)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.Text = task.Answer
	msg.ReplyToMessageID = task.MessageID
	msgs := utils.SplitMessageByMaxSize(task.Answer, 4000)
	for _, m := range msgs {
		msg.Text = m
		b.safeSend(msg)
	}
}

func (b *Bot) safeSendMsg(chatID int64, text string) {
	b.safeSend(tgbotapi.NewMessage(chatID, text))
}

func (b *Bot) safeSend(msg tgbotapi.MessageConfig) {
	if msg.Text == "" {
		return
	}
	_, err := b.tgBot.Send(msg)
	if err == nil {
		return
	}
	msg.ParseMode = ""
	_, err = b.tgBot.Send(msg)
	if err != nil {
		log.Printf("[Send] send message failed, err: 【%s】, msg: 【%+v】", err, msg)
		return
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[Update] 【chat】:%s, 【from】:%s, 【msg】:%s", utils.ToJsonString(update.Message.Chat),
		utils.ToJsonString(update.Message.From),
		utils.ToJsonString(update.Message))
	if update.Message.IsCommand() {
		b.execCommand(update.Message.Command(), update)
	} else {
		b.handleMessage(update)
	}

}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	log.Printf("[HandleMessage] [%s] update id[%d], from id[%d], from name[%s], msg[%s], chat id[%d], chat name[%s]",
		update.Message.Chat.Type, update.UpdateID,
		update.Message.From.ID, fmt.Sprintf("%s %s %s", update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName),
		update.Message.Text, update.Message.Chat.ID, update.Message.Chat.Title)

	_, _ = b.session.Load(update.Message.From.ID)

	if shouldIgnoreMsg(update) {
		return
	}
	ok := b.checkLimitersV1(*update.Message)
	if !ok {
		return
	}
	b.sendTaskToChannel(update.Message.Text, update.Message.Chat.ID, update.Message.From.ID, update.Message.MessageID)

	//if shouldHandleMessage(update, b.tgBot.Self.ID) {
	//
	//		b.sendTaskToChannel(update.Message.Text, update.Message.Chat.ID, update.Message.From.ID, update.Message.MessageID)
	//	} else {
	//		log.Printf("[RateLimit] user %d is chatting with me, ignore message %s", update.Message.From.ID, update.Message.Text)
	//		b.sendErrorMessage(update.Message.Chat.ID, update.Message.MessageID, constant.OnlyOneChatAtATime)
	//	}
	//}

}

func (b *Bot) sendErrorMessage(chatID int64, msgID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = msgID
	_, err := b.tgBot.Send(msg)
	if err != nil {
		log.Printf("[Send] send message failed, err: 【%s】, msg: 【%+v】", err, msg)
		msg.Text = constant.SendBackMsgFailed
		_, _ = b.tgBot.Send(msg)
		return
	}
}

func (b *Bot) findMemberFromChat(chatName string, userID int64) bool {
	findUserConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			SuperGroupUsername: chatName,
			UserID:             userID,
		},
	}
	member, err := b.tgBot.GetChatMember(findUserConfig)
	if err != nil || member.Status == "left" || member.Status == "kicked" {
		log.Printf("[ShouldLimitUser] memeber should be limit. id: %d", userID)
		return false
	}
	return true
}

func (b *Bot) shouldLimitUser(update tgbotapi.Update) bool {
	userID := strconv.FormatInt(update.Message.From.ID, 10)
	avaliable, err := b.userRepository.IsAvaliable(userID)
	if err != nil {
		log.Printf("[LimitUser] query user is avaliable failed, err: 【%s】\n", err)
		return false
	}
	return avaliable
}

func shouldIgnoreMsg(update tgbotapi.Update) bool {
	if update.Message == nil {
		return true
	}

	if update.Message.NewChatMembers != nil ||
		update.Message.LeftChatMember != nil {
		return true
	}

	if strings.Trim(update.Message.Text, " ") == "" {
		return true
	}

	return update.Message.ReplyToMessage != nil &&
		!update.Message.ReplyToMessage.From.IsBot
}

func (b *Bot) sendTaskToChannel(question string, chat, from int64, msgID int) {
	b.session.Store(from, &struct{}{})
	log.Printf("[SendTaskToChannel] with question %s, chat id: %d, from: %d", question, chat, from)
	chatTask := model.NewChatTask(question, chat, from, msgID)
	user, err := b.getUserInfo(from)
	if err == nil {
		chatTask.User = user
	}
	b.taskChan <- chatTask
	b.sendTyping(chatTask)
}

func (b *Bot) sendTyping(task *model.ChatTask) {
	msg := tgbotapi.NewChatAction(task.Chat, tgbotapi.ChatTyping)
	_, _ = b.tgBot.Send(msg)
}

//func (b *Bot) checkLimiters(update tgbotapi.Update) bool {
//	from := update.Message.From.ID
//	if update.Message.Chat.IsPrivate() {
//		if b.shouldLimitUser(update) {
//			b.sendErrorMessage(update.Message.Chat.ID, update.Message.MessageID, text)
//			return false
//		}
//	}
//	if b.enableLimiter &&
//		!b.limiter.Allow(strconv.FormatInt(update.Message.From.ID, 10)) {
//		log.Printf("[RateLimit] user %d is chatting with me, ignore message %s", update.Message.From.ID, update.Message.Text)
//		text := fmt.Sprintf(constant.RateLimitMessageTemplate,
//			b.limiter.GetCapacity(), b.limiter.GetDuration()/60,
//			b.channelName, b.groupName,
//			b.limiter.GetDuration()/60, b.limiter.GetCapacity(),
//			b.channelName, b.groupName)
//		b.sendErrorMessage(update.Message.Chat.ID, update.Message.MessageID, text)
//		return false
//	}
//	return true
//}

func (b *Bot) getUserInfo(userID int64) (*model.User, error) {
	user, err := b.tgBot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: userID,
		}})
	if err != nil {
		return nil, err
	}
	return model.NewUser(user.UserName, user.FirstName, user.LastName), nil
}

func (b *Bot) isBotAdmin(from int64) bool {
	if b.admin == 0 {
		return false
	}
	return b.admin == from
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
func (b *Bot) checkLimitersV1(m tgbotapi.Message) bool {
	for _, limiter := range b.limiters {
		ok, err := limiter.Allow(b, m)
		if !ok {
			if !utils.IsEmpty(err) {
				b.sendErrorMessage(m.Chat.ID, m.MessageID, err)
			}
			return false
		}
	}
	return true
}
