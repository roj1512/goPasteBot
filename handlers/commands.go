/**
 * goPasteBot - Telegram pastebin bot for https://ezup.dev/p/
 * Copyright (C) 2021  Roj Serbest
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package handlers

import (
	"bytes"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"goPasteBot/paster"
	"io"
	"net/http"
	"strings"
	"time"
)

const DELAY = 6 * time.Second

func deleteAfterDelay(b *gotgbot.Bot, m *gotgbot.Message) {
	time.Sleep(DELAY)
	m.Delete(b)
}

func ezPaste(b *gotgbot.Bot, m *gotgbot.Message) string {
	content, result := "", ""
	if m.Document != nil && m.Document.FileSize <= 9000 && strings.Split(m.Document.MimeType, "/")[0] == "text" {
		file, err := b.GetFile(m.Document.FileId)
		if err != nil {
			return ""
		}
		downloadUrl := "https://api.telegram.org/file/bot" + b.Token + "/" + file.FilePath
		resp, err := http.Get(downloadUrl)
		if err != nil {
			return ""
		}
		var buffer bytes.Buffer
		_, err = io.Copy(&buffer, resp.Body)
		if err != nil {
			return ""
		}
		content = buffer.String()
	} else if m.Text != "" {
		content = m.Text
	}
	if content == "" {
		return content
	}
	err := paster.Paste(content, &result)
	if err != nil {
		return ""
	}
	return result
}

func createShareUrl(url string) string {
	return "https://t.me/share/url?url=" + url +
		"&text=%E2%80%94%20__Pasted%20with__" +
		"%20%F0%9F%A4%96%20%40ezpastebot"
}

func paste(b *gotgbot.Bot, ctx *ext.Context) error {
	reply := ctx.EffectiveMessage.ReplyToMessage
	if reply.Text == "Send paste" {
		reply = ctx.EffectiveMessage
	}
	inputIsNotValid := !(reply != nil && (reply.Text != "" || reply.Document != nil))
	if inputIsNotValid {
		res, err := ctx.EffectiveMessage.Reply(
			b,
			"Reply to a text message/file with the command to "+
				"upload it to [ezpaste](https://ezup.dev/p/).",
			&gotgbot.SendMessageOpts{
				ParseMode:             "markdown",
				DisableWebPagePreview: true,
			},
		)
		if err == nil {
			deleteAfterDelay(b, res)
		}
	} else {
		url := ezPaste(b, reply)
		if url == "" {
			ctx.EffectiveMessage.Reply(b, "Invalid", nil)
		} else {
			shareUrl := createShareUrl(url)
			ctx.EffectiveMessage.Reply(
				b,
				url,
				&gotgbot.SendMessageOpts{
					ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
						InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
							{
								gotgbot.InlineKeyboardButton{
									Text: "Share",
									Url:  shareUrl,
								},
								gotgbot.InlineKeyboardButton{
									Text:              "Inline",
									SwitchInlineQuery: url,
								},
							},
						},
					},
				},
			)
		}
	}
	return nil
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	if strings.HasSuffix(ctx.EffectiveMessage.Text, "from_inline") {
		ctx.EffectiveMessage.Reply(
			b,
			ASK_TO_SEND_PASTE,
			&gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.ForceReply{ForceReply: true},
			},
		)
		return nil
	}
	ctx.EffectiveMessage.Reply(
		b,
		"🏷️ <b>How to use this bot to upload paste to "+
			"<a href=\"https://ezup.dev/p\">ezPaste</a></b> "+
			"(any of the following methods works):\n\n"+
			"- Use in inline mode\n"+
			"- send text or text file in private\n"+
			"- reply to a text message or text file with /paste in private "+
			"or groups (feel free to add this bot to your groups, it has "+
			"privacy mode enabled so it does not read your chat history\n\n"+
			"You can upload up to 1 megabytes of text on each paste\n\n"+
			"<a href=\"https://github.com/dashezup/ezpastebot)\">Source Code</a>"+
			" | <a href=\"https://t.me/dashezup\">Developer</a>"+
			" | <a href=\"https://t.me/ezupdev\">Support Chat</a>",
		&gotgbot.SendMessageOpts{
			DisableWebPagePreview: true,
			ParseMode:             "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						gotgbot.InlineKeyboardButton{
							Text:              "Try Inline Mode",
							SwitchInlineQuery: " ",
						},
					},
				},
			},
		},
	)
	return nil
}
