package bot

import (
	"chatgpt-bot/cfg"
	"chatgpt-bot/constant"
	"chatgpt-bot/engine"
	"chatgpt-bot/middleware"
	"chatgpt-bot/model"
	"chatgpt-bot/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	tgBot        *tgbotapi.BotAPI
	engine       engine.Engine
	session      *sync.Map
	taskChan     chan *model.ChatTask
	maxQueueChan chan interface{}
	limiter      *middleware.Limiter

	groupName     string
	channelName   string
	limitGroup    bool
	limitPrivate  bool
	logChat       int64
	enableLimiter bool
}

func (t *TelegramBot) Init(cfg *cfg.Config) error {

	t.channelName = cfg.BotConfig.TelegramChannelName
	t.groupName = cfg.BotConfig.TelegramGroupName
	t.limitPrivate = cfg.BotConfig.ShouldLimitPrivate
	t.limitGroup = cfg.BotConfig.ShouldLimitGroup
	t.logChat = cfg.BotConfig.LogChannelID

	if utils.IsAnyStringEmpty(
		t.channelName, t.groupName) {
		return errors.New(constant.MissingRequiredConfig)
	}

	t.session = &sync.Map{}
	bot, err := tgbotapi.NewBotAPI(cfg.BotConfig.TelegramBotToken)
	if err != nil {
		return err
	}
	t.tgBot = bot
	t.engine = engine.GetEngine(cfg.EngineConfig.EngineType)
	err = t.engine.Init(cfg)
	if err != nil {
		return err
	}

	t.taskChan = make(chan *model.ChatTask, 1)
	t.maxQueueChan = make(chan interface{}, 5)

	t.enableLimiter = cfg.BotConfig.RateLimiterConfig.Enable
	t.limiter = middleware.NewLimiter(cfg.BotConfig.RateLimiterConfig.Capacity,
		cfg.BotConfig.RateLimiterConfig.Duration)
	go t.loopAndFinishChatTask()
	log.Printf("[Init] telegram bot init success, bot name: %s", t.tgBot.Self.UserName)
	log.Printf("[Init] telegram bot init success, bot name: %s", t.tgBot.Self.UserName)
	return nil
}

func NewTelegramBot() *TelegramBot {
	return &TelegramBot{}
}

func (t *TelegramBot) Run() {
	log.Println("[Run] start telegram bot")
	go t.fetchUpdates()
}

func (t *TelegramBot) fetchUpdates() {
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60

	botChannel := t.tgBot.GetUpdatesChan(config)
	for {
		select {
		case update, ok := <-botChannel:
			if !ok {
				botChannel = t.tgBot.GetUpdatesChan(config)
				log.Println("[FetchUpdates] channel closed, fetch again")
				continue
			}
			go t.handleUpdate(update)
		case <-time.After(30 * time.Second):
		}
	}
}

func (t *TelegramBot) loopAndFinishChatTask() {
	for {
		select {
		case task := <-t.taskChan:
			log.Println("[LoopAndFinishChatTask] got a task to finish")
			go t.Finish(task)
		case <-time.After(30 * time.Second):
		}

	}
}

func (t *TelegramBot) Finish(task *model.ChatTask) {
	t.maxQueueChan <- struct{}{}
	defer func() {
		<-t.maxQueueChan
	}()
	log.Printf("[Finish] start chat task %s", task.String())
	defer t.session.Delete(task.From)
	t.Log(task.GetFormattedQuestion())
	t.sendTyping(task)
	res, err := t.engine.Chat(task.Question, strconv.FormatInt(task.From, 10))
	if err != nil {
		log.Printf("[Finish] chat task failed, err: %s", err)
		task.Answer = err.Error()
	} else {
		task.Answer = res
	}
	t.sendTyping(task)
	t.Send(task)
	t.Log(task.GetFormattedAnswer())
	log.Printf("[Finish] end chat task: %s", task.String())

}

func (t *TelegramBot) Log(log string) {
	go func(s string) {
		msg := tgbotapi.NewMessage(t.logChat, s)
		msg.ParseMode = tgbotapi.ModeMarkdown
		_ = t.safeSend(msg)
	}(log)
}

func (t *TelegramBot) Send(task *model.ChatTask) {
	msg := tgbotapi.NewMessage(task.Chat, task.Question)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.Text = task.Answer
	msg.ReplyToMessageID = task.MessageID
	msgs := utils.SplitMessageByMaxSize(task.Answer, 4000)
	for _, m := range msgs {
		msg.Text = m
		_ = t.safeSend(msg)
	}
}

func (t *TelegramBot) safeSend(msg tgbotapi.MessageConfig) error {
	if msg.Text == "" {
		return nil
	}
	_, err := t.tgBot.Send(msg)
	if err == nil {
		return nil
	}
	msg.ParseMode = ""
	_, err = t.tgBot.Send(msg)
	if err != nil {
		log.Printf("[Send] send message failed, err: 【%s】, msg: 【%+v】", err, msg)
		msg.Text = constant.SendBackMsgFailed
		_, _ = t.tgBot.Send(msg)
		return err
	}
	return nil
}

func (t *TelegramBot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[Update] 【chat】:%s, 【from】:%s, 【msg】:%s", utils.ToJsonString(update.Message.Chat),
		utils.ToJsonString(update.Message.From),
		utils.ToJsonString(update.Message))
	if update.Message.IsCommand() {
		msg := t.handleCommandMsg(update)
		_ = t.safeSend(msg)
	} else {
		t.handleUserMessage(update)
	}

}

func (t *TelegramBot) handleCommandMsg(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case constant.START:
		msg.Text = constant.BotStartTip
	case constant.CHATGPT:
		msg.Text = constant.BotStartTip
	case constant.PING:
		msg.Text = constant.BotPingTip
	case constant.Limiter:
		t.enableLimiter = utils.ParseBoolString(update.Message.CommandArguments())
		msg.Text = fmt.Sprintf("limiter status is %t now", t.enableLimiter)
	case constant.PPROF:
		if m, success := t.dumpProfile(update.Message.Chat.ID); !success {
			msg.Text = m
		}

	default:
		msg.Text = constant.UnknownCmdTip
	}
	return msg
}

func shouldHandleMessage(update tgbotapi.Update, selfID int64) bool {
	isPrivate := update.Message.Chat.IsPrivate()
	shouldHandleMessage := isPrivate ||
		(update.Message.ReplyToMessage != nil &&
			update.Message.ReplyToMessage.From.ID == selfID)
	return shouldHandleMessage
}

func (t *TelegramBot) handleUserMessage(update tgbotapi.Update) {
	log.Printf("[HandleMessage] [%s] update id[%d], from id[%d], from name[%s], msg[%s], chat id[%d], chat name[%s]",
		update.Message.Chat.Type, update.UpdateID,
		update.Message.From.ID, fmt.Sprintf("%s %s %s", update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName),
		update.Message.Text, update.Message.Chat.ID, update.Message.Chat.Title)

	_, thisUserHasMessage := t.session.Load(update.Message.From.ID)

	if shouldIgnoreMsg(update) {
		return
	}

	if shouldHandleMessage(update, t.tgBot.Self.ID) {
		if !t.checkLimiters(update) {
			return
		}
		if !thisUserHasMessage {
			t.sendTaskToChannel(update.Message.Text, update.Message.Chat.ID, update.Message.From.ID, update.Message.MessageID)
		} else {
			log.Printf("[RateLimit] user %d is chatting with me, ignore message %s", update.Message.From.ID, update.Message.Text)
			t.sendErrorMessage(update.Message.Chat.ID, update.Message.MessageID, constant.OnlyOneChatAtATime)
		}
	}

}

func (t *TelegramBot) sendErrorMessage(chatID int64, msgID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = msgID
	_, err := t.tgBot.Send(msg)
	if err != nil {
		log.Printf("[Send] send message failed, err: 【%s】, msg: 【%+v】", err, msg)
		msg.Text = constant.SendBackMsgFailed
		_, _ = t.tgBot.Send(msg)
		return
	}
}

func (t *TelegramBot) findMemberFromChat(chatName string, userID int64) bool {
	findUserConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			SuperGroupUsername: chatName,
			UserID:             userID,
		},
	}
	member, err := t.tgBot.GetChatMember(findUserConfig)
	if err != nil || member.Status == "left" || member.Status == "kicked" {
		log.Printf("[ShouldLimitUser] memeber should be limit. id: %d", userID)
		return false
	}
	return true
}

func (t *TelegramBot) shouldLimitUser(update tgbotapi.Update) bool {
	userID := update.Message.From.ID
	canFindInChannel := t.findMemberFromChat(t.channelName, userID)
	canFindInGroup := t.findMemberFromChat(t.groupName, userID)

	return !(canFindInChannel && canFindInGroup)
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

func (t *TelegramBot) sendTaskToChannel(question string, chat, from int64, msgID int) {
	t.session.Store(from, &struct{}{})
	log.Printf("[SendTaskToChannel] with question %s, chat id: %d, from: %d", question, chat, from)
	chatTask := model.NewChatTask(question, chat, from, msgID)
	user, err := t.getUserInfo(from)
	if err == nil {
		chatTask.User = user
	}
	t.taskChan <- chatTask
	t.sendTyping(chatTask)
}

func (t *TelegramBot) sendTyping(task *model.ChatTask) {
	msg := tgbotapi.NewChatAction(task.Chat, tgbotapi.ChatTyping)
	_, _ = t.tgBot.Send(msg)
}

func (t *TelegramBot) checkLimiters(update tgbotapi.Update) bool {
	if update.Message.Chat.IsPrivate() {
		if t.shouldLimitUser(update) {
			text := fmt.Sprintf(constant.LimitUserMessageTemplate, t.channelName, t.groupName, t.channelName, t.groupName)
			t.sendErrorMessage(update.Message.Chat.ID, update.Message.MessageID, text)
			return false
		}
	}
	if t.enableLimiter &&
		!t.limiter.Allow(strconv.FormatInt(update.Message.From.ID, 10)) &&
		t.shouldLimitUser(update) {
		log.Printf("[RateLimit] user %d is chatting with me, ignore message %s", update.Message.From.ID, update.Message.Text)
		text := fmt.Sprintf(constant.RateLimitMessageTemplate,
			t.limiter.GetCapacity(), t.limiter.GetDuration()/60,
			t.channelName, t.groupName,
			t.limiter.GetDuration()/60, t.limiter.GetCapacity(),
			t.channelName, t.groupName)
		t.sendErrorMessage(update.Message.Chat.ID, update.Message.MessageID, text)
		return false
	}
	return true
}

func (t *TelegramBot) getUserInfo(userID int64) (*model.User, error) {
	user, err := t.tgBot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: userID,
		}})
	if err != nil {
		return nil, err
	}
	return model.NewUser(user.UserName, user.FirstName, user.LastName), nil
}

func (t *TelegramBot) dumpProfile(chat int64) (string, bool) {
	fileName := fmt.Sprintf("%d.pprof", time.Now().Unix())
	filePath := os.TempDir() + string(os.PathSeparator) + fileName
	tmpFile, err := os.Create(filePath)
	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
		_ = os.Remove(filePath)
	}(tmpFile)

	if err != nil {
		log.Printf("[DumpProfile] create temp file failed, err: 【%s】", err)
		return err.Error(), false
	}

	err = pprof.WriteHeapProfile(tmpFile)
	if err != nil {
		log.Printf("[DumpProfile] create temp file failed, err: 【%s】", err)
		return err.Error(), false
	}

	err = t.sendFile(chat, tmpFile)
	if err != nil {
		log.Printf("[DumpProfile] send file failed, err: 【%s】", err)
		return err.Error(), false
	}
	return tmpFile.Name(), true
}

// send file to chat
func (t *TelegramBot) sendFile(chatID int64, file *os.File) error {
	fileMsg := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(file.Name()))
	_, err := t.tgBot.Send(fileMsg)
	if err != nil {
		log.Printf("[SendFile] send file failed, err: 【%s】", err)
		return err
	}
	return nil
}
