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
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func AddHandlers(dp *ext.Dispatcher) {
	dp.AddHandler(handlers.NewCommand("start", start))
	dp.AddHandler(handlers.NewCommand("paste", paste))
	dp.AddHandler(
		handlers.NewMessage(
			func(msg *gotgbot.Message) bool {
				if msg.Chat.Type != "private" {
					return false
				}
				if msg.ReplyToMessage != nil {
					if msg.ReplyToMessage.Text == "Send paste" {
						return true
					}
				}
				return false
			},
			paste,
		),
	)
	dp.AddHandler(
		handlers.NewMessage(
			func(msg *gotgbot.Message) bool {
				if msg.Chat.Type != "private" {
					return false
				}
				if msg.ReplyToMessage == nil {
					if msg.Text != "" || msg.Document != nil {
						return true
					}
				}
				return false
			},
			askToPaste,
		),
	)
	dp.AddHandler(handlers.NewCallback(filters.Equal("upload_paste"), uploadPaste))
	dp.AddHandler(handlers.NewCallback(filters.Equal("ignore_paste"), ignorePaste))
	dp.AddHandler(handlers.NewInlineQuery(nil, inline))
}
