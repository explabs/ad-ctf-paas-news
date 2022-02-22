package sender

import (
	"fmt"
	"github.com/explabs/ad-ctf-paas-exploits/service/storage"
	"log"
	"os"
)

func SendNews(round int) error {
	news, err := storage.FindNews(round)
	if err != nil {
		return err
	}
	if news == nil {
		return nil
	}
	token := os.Getenv("TELEGRAM_TOKEN")
	chatId := os.Getenv("TELEGRAM_CHAT_ID")
	fmt.Println(token, chatId)
	if token == "" || chatId == "" {
		return fmt.Errorf("telegram credentials is empty")
	}
	bot := TelegramBot{
		TelegramBotToken: token,
		ChatID:           chatId,
	}
	if news.Filename != "" {
		if err := bot.LoadMessage(news.Filename); err != nil {
			return err
		}
	} else if news.Text != "" {
		bot.Text = news.Text
		fmt.Println(bot.Text)
	}
	log.Println("sending message")
	if err := bot.SendMessage(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
