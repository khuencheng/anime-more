package bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

func MainHandler(m *tb.Message, b *tb.Bot) {
	log.Println("message: ", m)
	b.Send(m.Sender, "动画推荐bot")
}

func RecommendHandler(m *tb.Message, b *tb.Bot) {
	log.Println("message: ", m.Text, "payload: ", m.Payload)
	b.Send(m.Sender, "TODO ")
}
