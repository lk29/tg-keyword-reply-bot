package main

import (
	"flag"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"tg-keyword-reply-bot/common"
	"tg-keyword-reply-bot/db"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
)

var bot *api.BotAPI
var gcron *cron.Cron

var (
	debug       bool
	superUserId int
)

func main() {
	botToken := flag.String("t", "", "your bot Token")
	flag.IntVar(&superUserId, "s", 0, "super manager Id")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.Parse()
	token := db.Init(*botToken)
	gcron = cron.New()
	gcron.Start()
	//start working
	start(token)
}

func start(botToken string) {
	var err error
	bot, err = api.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = debug
	log.Printf("Authorized on account: %s  ID: %d", bot.Self.UserName, bot.Self.ID)

	u := api.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic("Can't get Updates")
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		go processUpdate(&update)
	}
}

/**
 * Separate processing for each update
 */
func processUpdate(update *api.Update) {
	upmsg := update.Message
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	//Check if it is a newly added group or a newly opened person
	in := checkInGroup(gid)
	if !in { //If it is not there, you need to add it, add a copy in the memory, and add an empty rule record in the database
		common.AddNewGroup(gid)
		db.AddNewGroup(gid)
	}
	if upmsg.IsCommand() {
		go processCommond(update)
	} else {
		go processReplyCommond(update)
		go processReply(update)
		//New users check if they are halal by username
		if upmsg.NewChatMembers != nil {
			for _, auser := range *(upmsg.NewChatMembers) {
				if checkQingzhen(auser.UserName) ||
					checkQingzhen(auser.FirstName) ||
					checkQingzhen(auser.LastName) {
					banMember(gid, uid, -1)
				}
			}
		}
		//Check halal and reject
		if checkQingzhen(upmsg.Text) {
			_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
			banMember(gid, uid, -1)
		}
	}
}

func processReply(update *api.Update) {
	var msg api.MessageConfig
	upmsg := update.Message
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	replyText := findKey(gid, upmsg.Text)
	if replyText == "delete" {
		_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
	} else if strings.HasPrefix(replyText, "ban") {
		_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
		banMember(gid, uid, -1)
	} else if strings.HasPrefix(replyText, "kick") {
		_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
		kickMember(gid, uid)
	} else if strings.HasPrefix(replyText, "photo:") {
		sendPhoto(gid, replyText[6:])
	} else if strings.HasPrefix(replyText, "gif:") {
		sendGif(gid, replyText[4:])
	} else if strings.HasPrefix(replyText, "video:") {
		sendVideo(gid, replyText[6:])
	} else if strings.HasPrefix(replyText, "file:") {
		sendFile(gid, replyText[5:])
	} else if replyText != "" {
		msg = api.NewMessage(gid, replyText)
		msg.DisableWebPagePreview = true
		msg.ReplyToMessageID = upmsg.MessageID
		sendMessage(msg)
	}
}

func processCommond(update *api.Update) {
	var msg api.MessageConfig
	upmsg := update.Message
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	msg = api.NewMessage(update.Message.Chat.ID, "")
	_, _ = bot.DeleteMessage(api.NewDeleteMessage(update.Message.Chat.ID, upmsg.MessageID))
	switch upmsg.Command() {
	case "start", "help":
		msg.Text = "This robot can automatically reply to specific keywords"
		sendMessage(msg)
	case "add":
		if checkAdmin(gid, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				addRule(gid, order)
				msg.Text = "The rule was added successfully: " + order
			} else {
				msg.Text = addText
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
			}
			sendMessage(msg)
		}
	case "del":
		if checkAdmin(gid, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				delRule(gid, order)
				msg.Text = "Rule deleted successfully: " + order
			} else {
				msg.Text = delText
				msg.ParseMode = "Markdown"
			}
			sendMessage(msg)
		}
	case "list":
		if checkAdmin(gid, *upmsg.From) {
			rulelists := getRuleList(gid)
			msg.Text = "ID: " + strconv.FormatInt(gid, 10)
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			sendMessage(msg)
			for _, rlist := range rulelists {
				msg.Text = rlist
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				sendMessage(msg)
			}
		}
	case "admin":
		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(uid) + ") Ask the admin to come out and spank\r\n\r\n" + getAdmins(gid)
		msg.ParseMode = "Markdown"
		sendMessage(msg)
		banMember(gid, uid, 30)
	case "banme":
		botme, _ := bot.GetChatMember(api.ChatConfigWithUser{ChatID: gid, UserID: bot.Self.ID})
		if botme.CanRestrictMembers {
			rand.Seed(time.Now().UnixNano())
			sec := rand.Intn(540) + 60
			banMember(gid, uid, int64(sec))
			msg.Text = "Congratulations[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ")get" + strconv.Itoa(sec) + "second silence package"
			msg.ParseMode = "Markdown"
		} else {
			msg.Text = "Please give me the ban permission, otherwise the game cannot be played"
		}
		sendMessage(msg)
	case "me":
		myuser := upmsg.From
		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ") account information for" +
			"\r\nID: " + strconv.Itoa(uid) +
			"\r\nUseName: [" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ")" +
			"\r\nLastName: " + myuser.LastName +
			"\r\nFirstName: " + myuser.FirstName +
			"\r\nIsBot: " + strconv.FormatBool(myuser.IsBot)
		msg.ParseMode = "Markdown"
		sendMessage(msg)
	default:
	}
}

func processReplyCommond(update *api.Update) {
	var msg api.MessageConfig
	upmsg := update.Message
	gid := upmsg.Chat.ID
	//Reply to administrative commands of type
	if upmsg.ReplyToMessage != nil {
		reolyToUserId := upmsg.ReplyToMessage.From.ID
		switch upmsg.Text {
		case "ban":
			if checkAdmin(gid, *upmsg.From) {
				banMember(gid, reolyToUserId, -1)
				mem, _ := bot.GetChatMember(api.ChatConfigWithUser{ChatID: gid, SuperGroupUsername: "", UserID: reolyToUserId})
				if !mem.CanSendMessages {
					msg = api.NewMessage(gid, "")
					msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ") muted " +
						"[" + upmsg.ReplyToMessage.From.String() + "](tg://user?id=" + strconv.Itoa(reolyToUserId) + ") "
					msg.ParseMode = "Markdown"
					sendMessage(msg)
				}
			}
		case "unban":
			if checkAdmin(gid, *upmsg.From) {
				unbanMember(gid, reolyToUserId)
				//mem,_ := bot.GetChatMember(api.ChatConfigWithUser{gid, "", reolyToUserId})
				//
				msg = api.NewMessage(gid, "")
				msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ") unlocked " +
					"[" + upmsg.ReplyToMessage.From.String() + "](tg://user?id=" + strconv.Itoa(reolyToUserId) + ") "
				msg.ParseMode = "Markdown"
				sendMessage(msg)
			}
		case "kick":
			if checkAdmin(gid, *upmsg.From) {
				kickMember(gid, reolyToUserId)
			}
		case "unkick":
			if checkAdmin(gid, *upmsg.From) {
				unkickMember(gid, reolyToUserId)
			}
		default:
		}
	}
}
