package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/satori/go.uuid"
	"strings"
)

const ASK_TO_SEND_PASTE = "Send paste"

func inline(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.InlineQuery.Query
	if strings.HasPrefix(query, "https://ezup.dev/p/") && len(query) == 25 {
		url := query
		rawUrl := query + "/index.txt"
		pasteId := url[10:]
		pasteInfo := "ezPaste: <a href=\"" + url + "\">" + pasteId + "</a> | <a href=\"" + rawUrl + "\">raw</a>"
		shareUrl := createShareUrl(url)
		thumbUrl := url + "/preview.png"
		ctx.InlineQuery.Answer(
			b,
			[]gotgbot.InlineQueryResult{
				gotgbot.InlineQueryResultArticle{
					Id:    uuid.NewV4().String(),
					Title: "Send URL of this paste",
					InputMessageContent: gotgbot.InputTextMessageContent{
						MessageText: pasteInfo,
						ParseMode: "HTML",
					},
					Url:      url,
					ThumbUrl: thumbUrl,
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
			},
			&gotgbot.AnswerInlineQueryOpts{
				CacheTime: 0,
			},
		)
	} else {
		ctx.InlineQuery.Answer(
			b,
			[]gotgbot.InlineQueryResult{},
			&gotgbot.AnswerInlineQueryOpts{
				CacheTime:         1,
				SwitchPmText:      "Send paste in PM",
				SwitchPmParameter: "from_inline",
			},
		)
	}
	return nil
}
