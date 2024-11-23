package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api: api,
	}, nil
}

func (b *Bot) Start() {
	log.Printf("Bot started: @%s", b.api.Self.UserName)

	// 配置更新
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	// 获取更新通道
	updates := b.api.GetUpdatesChan(updateConfig)

	// 处理更新
	for update := range updates {
		go b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		b.handleMessage(update.Message)
	case update.CallbackQuery != nil:
		b.handleCallback(update.CallbackQuery)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// 处理命令
	if message.IsCommand() {
		b.handleCommand(message)
		return
	}

	// 处理普通消息
	b.handleTextMessage(message)
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	switch message.Command() {
	case "start":
		msg.Text = "欢迎使用本机器人！\n输入 /help 查看帮助信息。"

	case "help":
		msg.Text = `可用命令：
/start - 开始使用
/help - 显示帮助
/about - 关于
/keyboard - 显示键盘
/photo - 发送图片示例`

	case "about":
		msg.Text = "这是一个示例机器人，使用Golang开发。"

	case "keyboard":
		msg.ReplyMarkup = b.getKeyboard()
		msg.Text = "这是一个自定义键盘示例"

	case "photo":
		// 发送图片示例
		photo := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL("https://example.com/photo.jpg"))
		b.api.Send(photo)
		return

	default:
		msg.Text = "未知命令。输入 /help 查看可用命令。"
	}

	b.api.Send(msg)
}

func (b *Bot) handleTextMessage(message *tgbotapi.Message) {
	// 创建回复消息
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	// 处理不同的文本输入
	switch strings.ToLower(message.Text) {
	case "hello", "hi":
		msg.Text = fmt.Sprintf("你好，%s!", message.From.FirstName)

	case "时间":
		msg.Text = fmt.Sprintf("当前时间：%s", time.Now().Format("2006-01-02 15:04:05"))

	default:
		// 回显用户消息
		msg.Text = "你说: " + message.Text
	}

	b.api.Send(msg)
}

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) {
	// 处理按钮回调
	callback_response := tgbotapi.NewCallback(callback.ID, callback.Data)
	b.api.Request(callback_response)

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "")
	msg.Text = "你点击了: " + callback.Data

	b.api.Send(msg)
}

func (b *Bot) getKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Hello"),
			tgbotapi.NewKeyboardButton("时间"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("帮助"),
			tgbotapi.NewKeyboardButton("关于"),
		),
	)
	keyboard.ResizeKeyboard = true
	return keyboard
}

// 创建内联键盘
func (b *Bot) getInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("访问网站", "https://example.com"),
			tgbotapi.NewInlineKeyboardButtonData("点击我", "button1"),
		),
	)
}

func main() {
	// 从环境变量获取 token
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	token = "8108318198:AAHwP0xMe-281WdhGUdj9KgbDQspLq5XqG8"
	if token == "" {
		log.Fatal("请设置 TELEGRAM_BOT_TOKEN 环境变量")
	}

	// 创建机器人实例
	bot, err := NewBot(token)
	if err != nil {
		log.Fatal(err)
	}

	// 启动机器人
	bot.Start()
}
