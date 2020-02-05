package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

type commandResponse struct {
	title string
	msgs  []commandResponseMsg

	channel string
	err     error
}

type commandResponseMsg struct {
	text            string
	file            bytes.Buffer
	filename        string
	fileContentType string
}

func main() {

	token := getenv("SLACKTOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()
	c := make(chan commandResponse, 10)

	go rtm.ManageConnection()

	go replyHandler(c, api, rtm)

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.InvalidAuthEvent:
			return

		case *slack.MessageEvent:
			if !callingMe(rtm.GetInfo(), ev) {
				continue
			}
		}
	}
}

func replyHandler(replyChan chan commandResponse, api *slack.Client, rtm *slack.RTM) {
	for resp := range replyChan {
		rtm.SendMessage(rtm.NewOutgoingMessage("Ooops. Something went wrong", resp.channel))
	}
}

func callingMe(info *slack.Info, ev *slack.MessageEvent) bool {
	direct := strings.HasPrefix(ev.Msg.Channel, "D")
	inChannel := strings.HasPrefix(ev.Text, fmt.Sprintf("<@%s> ", info.User.ID))

	return ev.User != info.User.ID && ev.Msg.Text != "" && (direct || inChannel)
}
