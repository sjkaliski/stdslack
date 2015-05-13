package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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

	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Check if stdin is from a terminal (i.e. not input to read).
	if stat.Mode()&os.ModeCharDevice != 0 {
		fmt.Fprintln(os.Stderr, "Content needs to be given to stdin to use")
		os.Exit(1)
	}

	var content bytes.Buffer
	_, err = io.Copy(&content, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = slackC.SendMessage(*channel, content.String(), "stdslack")
	if err != nil {
		fmt.Println(err)
	}
}
