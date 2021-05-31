package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func askToPaste(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.EffectiveMessage.Reply(
		b,
		"Do you want to upload this paste?",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						gotgbot.InlineKeyboardButton{
							Text:         "Yes",
							CallbackData: "upload_paste",
						},
						gotgbot.InlineKeyboardButton{
							Text:         "No",
							CallbackData: "ignore_paste",
						},
					},
				},
			},
		},
	)
	return nil
}

func uploadPaste(b *gotgbot.Bot, ctx *ext.Context) error {
	url := ezPaste(b, ctx.EffectiveMessage.ReplyToMessage)
	if url == "" {
		return nil
	}
	shareUrl := createShareUrl(url)
	ctx.EffectiveMessage.EditText(
		b,
		url,
		&gotgbot.EditMessageTextOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						gotgbot.InlineKeyboardButton{
							Text: "Share",
							Url:  shareUrl,
						},
					},
				},
			},
		},
	)
	return nil
}

func ignorePaste(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.EffectiveMessage.Delete(b)
	return nil
}
