package main

import (
	"anime-more/backend"
	"anime-more/bot"
	"anime-more/config"
)

func main() {
	conf := config.GetConfig()
	bot.StartBot()
	backend.StartServer(true, conf.GetString("web.http_port"))
}
