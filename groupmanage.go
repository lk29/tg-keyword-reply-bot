package main

import (
	"strconv"
	"tg-keyword-reply-bot/common"
	"time"
	"unicode"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

/**
 * Check if you are an admin of the group
 */
func checkAdmin(gid int64, user api.User) bool {
	admins, _ := bot.GetChatAdministrators(api.ChatConfig{ChatID: gid, SuperGroupUsername: ""})
	uid := user.ID
	if superUserId > 0 && uid == superUserId {
		return true
	}
	for _, user := range admins {
		if uid == user.User.ID {
			return true
		}
	}
	chengfa(gid, user)
	return false
}

/**
 * Check if it is a newly added group or a newly opened person
 */
func checkInGroup(id int64) bool {
	for _, gid := range common.AllGroupId {
		if gid == id {
			return true
		}
	}
	return false
}

func chengfa(gid int64, user api.User) {
	botme, _ := bot.GetChatMember(api.ChatConfigWithUser{ChatID: gid, UserID: bot.Self.ID})
	msg := api.NewMessage(gid, "")
	if botme.CanRestrictMembers {
		banMember(gid, user.ID, 60)
		msg.Text = "[" + user.String() + "](tg://user?id=" + strconv.Itoa(user.ID) + ")Randomly playing with administrator commands, ban for one minute"
		msg.ParseMode = "Markdown"
	} else {
		msg.Text = "[" + user.String() + "](tg://user?id=" + strconv.Itoa(user.ID) + ")Don't mess around with admin commands"
		msg.ParseMode = "Markdown"
	}
	sendMessage(msg)
}

/**
 * Forbidden group members
 */
func banMember(gid int64, uid int, sec int64) {
	if sec <= 0 {
		sec = 9999999999999
	}
	chatuserconfig := api.ChatMemberConfig{ChatID: gid, UserID: uid}
	b := false
	restricconfig := api.RestrictChatMemberConfig{
		ChatMemberConfig:      chatuserconfig,
		UntilDate:             time.Now().Unix() + sec,
		CanSendMessages:       &b,
		CanSendMediaMessages:  &b,
		CanSendOtherMessages:  &b,
		CanAddWebPagePreviews: &b}
	_, _ = bot.RestrictChatMember(restricconfig)
}

func unbanMember(gid int64, uid int) {
	chatuserconfig := api.ChatMemberConfig{ChatID: gid, UserID: uid}
	b := true
	restricconfig := api.RestrictChatMemberConfig{
		ChatMemberConfig:      chatuserconfig,
		UntilDate:             9999999999999,
		CanSendMessages:       &b,
		CanSendMediaMessages:  &b,
		CanSendOtherMessages:  &b,
		CanAddWebPagePreviews: &b}
	_, _ = bot.RestrictChatMember(restricconfig)
}

/**
 * Kick out group members
 */
func kickMember(gid int64, uid int) {
	cmconf := api.ChatMemberConfig{ChatID: gid, UserID: uid}
	_, _ = bot.KickChatMember(api.KickChatMemberConfig{ChatMemberConfig: cmconf, UntilDate: 99999999999})
}

func unkickMember(gid int64, uid int) {
	_, _ = bot.UnbanChatMember(api.ChatMemberConfig{ChatID: gid, UserID: uid})
}

/**
 * Returns all admins of the group, for one-time @
 */
func getAdmins(gid int64) string {
	admins, _ := bot.GetChatAdministrators(api.ChatConfig{ChatID: gid})
	list := ""
	for _, admin := range admins {
		user := admin.User
		if user.IsBot {
			continue
		}
		list += "[" + user.String() + "](tg://user?id=" + strconv.Itoa(admin.User.ID) + ")\r\n"
	}
	return list
}

/**
 * Check if the text contains Arabic
 */
func checkQingzhen(text string) bool {
	for _, c := range text {
		if unicode.Is(unicode.Scripts["Arabic"], c) {
			return true
		}
	}
	return false
}
