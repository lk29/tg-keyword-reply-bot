package main

import (
	"log"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

/**
* Send text message 
*/
func sendMessage(msg api.MessageConfig) api.Message {
	if msg.Text == "" {
		return api.Message{}
	}
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	go deleteMessage(msg.ChatID, mmsg.MessageID)
	return mmsg
}

/**
 * To send a picture message, it needs to be an existing picture link
 */
func sendPhoto(chatid int64, photoid string) api.Message {
	file := api.NewPhotoShare(chatid, photoid)
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Println(err)
	}
	go deleteMessage(chatid, mmsg.MessageID)
	return mmsg
}

/**
 * To send an animation, it needs to be an existing link
 */
func sendGif(chatid int64, gifid string) api.Message {
	file := api.NewAnimationShare(chatid, gifid)
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Println(err)
	}
	go deleteMessage(chatid, mmsg.MessageID)
	return mmsg
}

/**
 * To send video, it needs to be an existing video connection
 */
func sendVideo(chatid int64, videoid string) api.Message {
	file := api.NewVideoShare(chatid, videoid)
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Println(err)
	}
	go deleteMessage(chatid, mmsg.MessageID)
	return mmsg
}

/**
 * To send a file, it must be an existing file link
 */
func sendFile(chatid int64, fileid string) api.Message {
	file := api.NewDocumentShare(chatid, fileid)
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Println(err)
	}
	go deleteMessage(chatid, mmsg.MessageID)
	return mmsg
}

func deleteMessage(gid int64, mid int) {
	time.Sleep(time.Second * 240)
	_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, mid))
}
