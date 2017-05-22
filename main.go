package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/nlopes/slack"
)

var (
	enableSlackBot = flag.Bool("slackBot", false, "Enable slackbot integration with service")
)

func main() {
	flag.Parse()

	if *enableSlackBot {
		commands := parseCommandsJSON("./commands.json")
		go initSlackBot(commands)
	}

	// start web server for handling Github webhook
	http.HandleFunc("/webhook", webhookHandler)
	http.ListenAndServe(":8080", nil)
}

/*
 * In webhook we want to parse the JSON request body in order to run testing.
 * from the webhook request, we will want to determine student public URL (this
 * may require a temporary storage to store student repo to its public URL
 * relationship). From the public URL, we will trigger the WebDriverIO test
 * against it. Upon the finish of the test, store the result and its SHA for
 * future reference to see its build status.
 * Last but not least, we will be posting the status back to Github to complete
 * the entire webhook event.
 */
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// parse request body
	// Get student public URL
	// run Webdriver IO test against the URL and get its result
	// Store result and SHA
	// publish status to Github
	fmt.Println("Hello")
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
