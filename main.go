package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	// A mapping between command name and its shell command
	var commands map[string]string
	// read commands.json to commands hashmap
	// TODO: maybe change this commands.json to command line argument
	file, err := ioutil.ReadFile("./commands.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(file, &commands)
	log.Printf("Got commands:\n%q\n", commands)

	api := slack.New(os.Getenv("SLACK_API_KEY"))
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)
	var selfID string

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			log.Println("Connected Infos:", ev.Info)
			selfID = ev.Info.User.ID
			log.Println("User ID: ", selfID)
			rtm.SendMessage(rtm.NewOutgoingMessage("TA bot is online", "C55V47YU9"))

		case *slack.MessageEvent:
			log.Printf("Message: %v\n", ev)
			// skip self talk
			if ev.User == selfID {
				break
			}
			// only reply to direct message
			if strings.HasPrefix(ev.Channel, "D") {
				result, err := handle(commands, ev.Text)
				if err != nil {
					result = "Error while executing command. Plesae let human know."
					log.Fatal(err)
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(
					result,
					ev.Channel,
				))
			}

		case *slack.PresenceChangeEvent:
			// log.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			// log.Printf("Current latency: %v\n", ev.Value)

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

func handle(commands map[string]string, text string) (string, error) {
	log.Printf("Got message content: %v\n", text)
	commandParts := strings.Split(commands[text], " ")
	var output bytes.Buffer
	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	cmd.Stdout = &output
	log.Printf("Running command: `%q`", commandParts)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
