package crawler

import (
	"anime-more/config"
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func init() {
	config.Init("dev")
	conf := config.GetConfig()
	soup.Header("User-Agent", conf.GetString("downloader.useragent"))
}

type Item struct {
	Title string
	Url   string
	Pic   string
}

func getDoubanPageLink(keyword string) string {
	link := config.GetConfig().GetString("urls.douban_search") + keyword
	resp, err := soup.Get(link)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	doc := soup.HTMLParse(resp)
	element := doc.Find("div", "class", "result")
	a := element.Find("a", "class", "nbg")
	txt := a.Attrs()["onclick"]
	r, _ := regexp.Compile(`sid: (\d+)`)
	txt = r.FindString(txt)
	txt = strings.Replace(txt, `sid: `, "", 1)
	return config.GetConfig().GetString("urls.douban") + txt
}

func DownloadDouban(keyword string) []Item {
	link := getDoubanPageLink(keyword)
	resp, err := soup.Get(link)
	if err != nil {
		fmt.Println(err)
		return []Item{}
	}
	doc := soup.HTMLParse(resp)
	elements := doc.Find("div", "class", "recommendations-bd").FindAll("a")
	if len(elements) == 0 {
		return []Item{}
	}
	items := make([]Item, 0, 20)
	for _, element := range elements {
		if element.Text() == "" {
			img := element.Find("img")
			items = append(items, Item{
				Title: img.Attrs()["alt"],
				Url:   element.Attrs()["href"],
				Pic:   img.Attrs()["src"],
			})
		}
	}
	return items
}

func parseBilibiliSeasonID(txt string) string {
	r, _ := regexp.Compile(`"season_id":(\d+)`)
	found := r.FindString(txt)
	found = strings.Replace(found, `"season_id":`, "", 1)
	return found
}

func DownloadBiliBili(keyword string) []Item {
	link := config.GetConfig().GetString("urls.bilibili") + keyword
	fmt.Println(link)
	resp, err := soup.Get(link)
	if err != nil {
		fmt.Println(err)
		return []Item{}
	}
	doc := soup.HTMLParse(resp)
	seasonID := parseBilibiliSeasonID(doc.FullText())
	recommendUrl := config.GetConfig().GetString("urls.bilibili_recommend") + seasonID
	fmt.Println(recommendUrl)
	type Recommend struct {
		Code int `json:"code"`
		Data struct {
			Season []struct {
				Actor      string `json:"actor"`
				Cover      string `json:"cover"`
				SeasonID   int    `json:"season_id"`
				SeasonType int    `json:"season_type"`
				Subtitle   string `json:"subtitle"`
				Title      string `json:"title"`
				URL        string `json:"url"`
			} `json:"season"`
		} `json:"data"`
	}
	res, err := http.Get(recommendUrl)
	if err != nil {
		fmt.Println(err)
		return []Item{}
	}
	defer res.Body.Close()
	recommendData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []Item{}
	}
	var recommend Recommend
	err = json.Unmarshal(recommendData, &recommend)
	if err != nil {
		fmt.Println(err)
		return []Item{}
	}
	items := make([]Item, 0, 20)
	for _, rec := range recommend.Data.Season {
		items = append(items, Item{
			Title: rec.Title,
			Url:   rec.URL,
			Pic:   rec.Cover,
		})
	}
	return items

}

func getMALRecommendLink(keyword string) string {
	link := config.GetConfig().GetString("urls.myanimelist_search") + keyword
	resp, err := soup.Get(link)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	doc := soup.HTMLParse(resp)
	element := doc.Find("div", "class", "js-categories-seasonal").Find("a", "class", "hoverinfo_trigger")
	link = element.Attrs()["href"] + "/userrecs"
	return link
}

func DownloadMAL(keyword string) []Item {
	link := getMALRecommendLink(keyword)
	resp, err := soup.Get(link)
	if err != nil {
		fmt.Println(err)
		return []Item{}
	}
	doc := soup.HTMLParse(resp)
	elements := doc.Find("div", "id", "contentWrapper").FindAll("div", "class", "picSurround")
	if len(elements) == 0 {
		return []Item{}
	}
	items := make([]Item, 0, 3)
	for i, element := range elements {
		imgTag := element.Find("img")
		a := element.Find("a")
		if i < 30 {
			items = append(items, Item{
				Title: strings.Replace(imgTag.Attrs()["alt"], "Anime: ", "", 1),
				Url:   a.Attrs()["href"],
				Pic:   imgTag.Attrs()["data-src"],
			})
		}

	}
	return items

}

func DownloadIBangumi(link string) []Item {
	resp, err := soup.Get(link)
	if err != nil {
		fmt.Println(err)
		return []Item{}
	}
	doc := soup.HTMLParse(resp)
	elements := doc.Find("ul", "class", "coversSmall").FindAll("a", "class", "thumbTip")
	if len(elements) == 0 {
		return []Item{}
	}
	items := make([]Item, 0, 20)
	for i, element := range elements {
		if i < 50 {
			picUrl := element.Find("span").Attrs()["style"]
			picUrl = strings.ReplaceAll(picUrl, `background-image:url('`, "https:")
			picUrl = strings.ReplaceAll(picUrl, "')", "")
			items = append(items, Item{
				Title: element.Attrs()["title"],
				Url:   "https://bgm.tv/" + element.Attrs()["href"],
				Pic:   picUrl,
			})
		}

	}
	return items

}

func DownloadIMDB(link string) []Item {
	resp, err := soup.Get(link)
	if err != nil {
		fmt.Println(err)
		return []Item{}
	}
	doc := soup.HTMLParse(resp)
	elements := doc.Find("div", "class", "rec_view").FindAll("div", "class", "rec_item")
	if len(elements) == 0 {
		return []Item{}
	}
	items := make([]Item, 0, 50)
	for i, element := range elements {
		a := element.Find("a")
		imgTag := a.Find("img")
		if i < 50 {
			items = append(items, Item{
				Title: imgTag.Attrs()["title"],
				Url:   "https://www.imdb.com/" + a.Attrs()["href"],
				Pic:   imgTag.Attrs()["src"],
			})
		}

	}
	return items

}
