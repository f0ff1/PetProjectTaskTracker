package handler

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

)

const (
	StateIdle                = ""
	StateAwaitingTitle       = "awaiting_title"
	StateAwaitingDescription = "awaiting_description"
	StateAwaitingTags        = "awaiting_tags"

	dateFormat = "02.01.2006 15:04"
)

type PendingTask struct {
	Title       string
	Description string
	Tags        []string
}


func (h *TelegramHandler) sendMessage(chatID int64, text string, markdown bool) {
	msg := tgbotapi.NewMessage(chatID, text)
	if markdown {
		msg.ParseMode = "Markdown"
	}
	h.bot.Send(msg)
}


func (h *TelegramHandler) sendMarkdown(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	h.bot.Send(msg)
}


func (h *TelegramHandler) parseTitleAndTags(raw string) (string, []string) {
	words := strings.Fields(raw)
	titleWords := []string{}
	tags := []string{}

	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#")
			if tag != "" {
				tags = append(tags, tag)
			}
		} else {
			titleWords = append(titleWords, word)
		}
	}

	title := strings.Join(titleWords, " ")
	return title, tags
}


func (h *TelegramHandler) parseTags(input string) []string {
	words := strings.Fields(input)
	tags := []string{}
	seen := make(map[string]bool)

	for _, word := range words {
		tag := strings.TrimPrefix(word, "#")
		tag = strings.TrimSpace(tag)
		if tag != "" && !seen[tag] {
			seen[tag] = true
			tags = append(tags, tag)
		}
	}
	return tags
}

// escapeMarkdown экранирует специальные символы для Markdown
func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}
