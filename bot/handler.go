package bot

import (
	"anime-more/crawler"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/url"
	"sync"
)

func MainHandler(m *tb.Message, b *tb.Bot) {
	log.Println("message: ", m)
	b.Send(m.Sender, "动画推荐bot")
}

func RecommendHandler(m *tb.Message, b *tb.Bot) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			b.Send(m.Sender, fmt.Sprintln("Error: ", err))
		}
	}()

	log.Println("message: ", m.Text, "payload: ", m.Payload)
	keyword := url.QueryEscape(m.Payload)
	items := make([]crawler.Item, 0, 100)
	var douban, bili, mal []crawler.Item
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		douban = crawler.DownloadDouban(keyword)
		wg.Done()
	}()
	go func() {
		bili = crawler.DownloadBiliBili(keyword)
		wg.Done()
	}()
	go func() {
		mal = crawler.DownloadMAL(keyword)
		wg.Done()

	}()
	wg.Wait()
	items = append(items, douban...)
	items = append(items, bili...)
	items = append(items, mal...)
	for _, item := range items {
		formated := fmt.Sprintf("<a href=\"%s\"> </a> \n \n <a href=\"%s\"><b>%s</b></a> \n \n 来自：%s", item.Pic, item.Url, item.Title, item.From)
		b.Send(m.Sender, formated, "HTML")
	}

}
