package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	// A mapping between command name and its shell command
	commands := parseCommandsJSON("./commands.json")
	go initSlackBot(commands)

	// start web server for handling Github webhook
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
}

func parseCommandsJSON(JSONFilePath string) map[string]string {
	var commands map[string]string
	file, err := ioutil.ReadFile(JSONFilePath)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(file, &commands)
	return commands
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
}

func initSlackBot(commands map[string]string) {
	var selfID string
	api := slack.New(os.Getenv("SLACK_API_KEY"))
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
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
				result, err := handleCommand(commands, ev.Text)
				if err != nil {
					result = err.Error()
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(
					result,
					ev.Channel,
				))
			}

		case *slack.InvalidAuthEvent:
			log.Printf("Invalid credentials")
			return

		default:
		}
	}
}

func handleCommand(commands map[string]string, text string) (string, error) {
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
