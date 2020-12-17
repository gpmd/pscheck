package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/keybase/go-ps"
	"github.com/nlopes/slack"
	"gopkg.in/yaml.v2"
)

func main() {
	conf := struct {
		Token  string `yaml:"token"`
		ChatID string `yaml:"chatid"`
	}{}

	data, err := ioutil.ReadFile(".pscheck")
	if err != nil {
		fmt.Println("Can't read .pscheck file")
		os.Exit(0)
	}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		fmt.Printf("Can't parse .pscheck file: %v\n", err)
		os.Exit(0)
	}
	if len(os.Args) < 3 {
		fmt.Println("Sends slack message is process with binary doesn't run.")
		fmt.Println("Usage: pscheck {executable name (example: lim)} {servername for chat message}")
		os.Exit(0)
	}
	pp, err := ps.Processes()
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	if len(pp) <= 0 {
		log.Fatalf("should have processes")
	}

	for _, p := range pp {
		if p.Executable() == os.Args[1] {
			os.Exit(0)
		}
	}
	api := slack.New(conf.Token)
	api.PostMessage(
		conf.ChatID,
		slack.MsgOptionText(fmt.Sprintf("PSCheck (from cron): No '%s' process running on server %s", os.Args[1], os.Args[2]), false),
	)
}
