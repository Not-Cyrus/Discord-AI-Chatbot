package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valyala/fastjson"

	"github.com/bwmarrin/discordgo"

	"github.com/valyala/fasthttp"
)

// r5fawAJCjWtSBfkF if you're here reading the code have this. (It'll probably be banned or something by the time you read this code lol)

func main() {
	fmt.Print("Enter your bot token: ")
	fmt.Scan(&token)
	fmt.Print("\nOK, now enter your brainshop.ai key (https://brainshop.ai): ")
	fmt.Scan(&apikey)
	ds, err := discordgo.New("Bot " + token)
	if err != nil {
		panic("wrong token")
	}
	ds.AddHandler(ready)
	ds.AddHandler(messageCreate)

	ds.Open()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	ds.Close()
}

func ready(s *discordgo.Session, ready *discordgo.Ready) {
	fmt.Println("We are ready.")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return // LOL I almost released it without this here
	}
	responseMessage := sendRequest("GET", fmt.Sprintf("http://api.brainshop.ai/get?bid=154371&key=%s&uid=%s&msg=%s", apikey, url.QueryEscape(m.Author.Username), url.QueryEscape(m.Content)))
	parsed, err := parser.Parse(string(responseMessage.Body()))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Couldn't parse json data: %s", err.Error()))
		return
	}
	s.ChannelMessageSend(m.ChannelID, string(parsed.GetStringBytes("cnt")))
}

func sendRequest(method, url string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.Header.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	err := client.DoTimeout(req, resp, 10*time.Second)
	if err != nil {
		fmt.Printf("Error doing the http request: %s\n", err.Error())
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		return nil
	}
	return resp
}

var (
	token  string
	apikey string
	parser fastjson.Parser
	client = &fasthttp.Client{}
)
