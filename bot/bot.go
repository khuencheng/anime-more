package bot

import (
	"anime-more/config"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

func StartBot() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	b, err := tb.NewBot(tb.Settings{
		Token:  config.GetConfig().GetString("telegram.bot_token"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	log.Println("bot", b)
	if err != nil {
		log.Println(err)
		return
	}

	b.Handle("/anime", func(m *tb.Message) {
		RecommendHandler(m, b)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		MainHandler(m, b)
	})

	b.Start()
}
