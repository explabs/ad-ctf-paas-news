package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const NewsFolder = "news"
const ImageFolder = "news/image"
const audioFolder = "news/audio"

type TelegramBot struct {
	TelegramBotToken string
	ChatID           string
	Text             string
}

type TelegramMessage struct {
	ChatID              string `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
	ParseMode           string `json:"parse_mode"`
}

func escapeMarkdown(text string) string {
	characters := map[string]string{
		"!": "\\!",
		".": "\\.",
		"{": "\\{",
		"}": "\\}",
		">": "\\>",
		"-": "\\-",
		"#": "\\#",
	}
	for old, _new := range characters {
		text = strings.Replace(text, old, _new, -1)
	}

	runes := []rune(text)
	for i := 1; i < len(runes); i++ {
		if runes[i-1] == '_' && runes[i] == '_' {

		}
	}

	return text
}

func (bot *TelegramBot) SendMessage() error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", bot.TelegramBotToken)
	body := TelegramMessage{
		ChatID:              bot.ChatID,
		Text:                bot.Text,
		DisableNotification: false,
		ParseMode:           "MarkdownV2",
	}

	jsonBody, _ := json.Marshal(body)
	fmt.Println(string(jsonBody))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	data, err := client.Do(req)
	if err != nil {
		return err
	}
	bodyBytes, err := io.ReadAll(data.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Info(bodyString)
	return nil
}

func (bot *TelegramBot) LoadMessage(fileName string) error {
	file, err := os.Open(filepath.Join(NewsFolder, fileName))
	if err != nil {
		return err
	}

	text, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	bot.Text = string(text)
	if err = file.Close(); err != nil {
		return err
	}
	return nil
}
