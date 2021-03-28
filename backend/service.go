package backend

import (
	"anime-more/crawler"
	"net/url"
	"sync"
)

func GetRecommendService(keyword string) []crawler.Item {
	keyword = url.QueryEscape(keyword)
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
	return items
}
