package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Bowery/slack"
)

var (
	channel    = flag.String("channel", "", "Channel to post to.")
	token      = flag.String("token", "", "Slack auth token.")
	configPath string
	slackToken string
	err        error
)

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	configPath = filepath.Join(u.HomeDir, ".stdslackconf")
}

func main() {
	flag.Parse()
	if *token != "" {
		err = ioutil.WriteFile(configPath, []byte(*token), 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(os.Stdout, "Wrote token to %s\n", configPath)
		return
	}

	if *channel == "" {
		fmt.Println("channel required.")
		return
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("run `stdslack --token=YOUR_TOKEN` to set token before using.")
		return
	}
	slackC := slack.NewClient(string(data))

	var content string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	err = slackC.SendMessage(*channel, content, "stdslack")
	if err != nil {
		fmt.Println(err)
	}
}
