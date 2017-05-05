package main

import (
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New(os.Getenv("SLACK_API_KEY"))
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			log.Println("Infos:", ev.Info)
			log.Println("Connection counter:", ev.ConnectionCount)
			rtm.SendMessage(rtm.NewOutgoingMessage("TA bot is online", "C55V47YU9"))

		case *slack.MessageEvent:
			log.Printf("Message: %v\n", ev)
			// only reply to direct message
			if strings.HasPrefix(ev.Channel, "D") {
				rtm.SendMessage(rtm.NewOutgoingMessage("pong", ev.Channel))
				handle(ev.Text)
			}

		case *slack.PresenceChangeEvent:
			log.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			log.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			log.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func handle(text string) {
	log.Printf("Got message content: %v\n", text)
}
