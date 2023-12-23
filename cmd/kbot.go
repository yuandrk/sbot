/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

var (
	TeleToken = os.Getenv("TELE_TOKEN")
)

var sbotCmd = &cobra.Command{
	Use:     "sbot",
	Aliases: []string{"start"},
	Short:   "Start telegram sbot application",
	Long: `A simple Telegram bot that can handle text messages.
	You can write something to https://t.me/yuandrk_bot and sometimes it answer:)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sbot %s started\n", appVersion)

		sbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("please check TELE_TOKEN env variable, %s", err)
		}

		sbot.Handle(telebot.OnText, func(ctx telebot.Context) error {
			log.Println(ctx.Message().Payload, ctx.Text())
			payload := ctx.Message().Payload
			answerStr := handlePayload(payload)
			err = ctx.Send(answerStr)
			return err
		})

		sbot.Handle(telebot.OnVoice, func(ctx telebot.Context) error {
			answerStr := "I don' have ears, so I can't hear you :)"
			err = ctx.Send(answerStr)
			return err
		})

		sbot.Handle(telebot.OnPhoto, func(ctx telebot.Context) error {
			answerStr := "Nice picture... or not"
			err = ctx.Send(answerStr)
			return err
		})

		sbot.Handle(telebot.OnSticker, func(ctx telebot.Context) error {
			answerStr := "It's very funny, i think so"
			err = ctx.Send(answerStr)
			return err
		})

		sbot.Start()
	},
}

func handlePayload(payload string) string {
	switch strings.ToLower(payload) {
	case "hello":
		return fmt.Sprintf("Hello I'm sbot %s", appVersion)
	case "version":
		return fmt.Sprintf("Current sbot version %s", appVersion)
	case "how are you?":
		return "Thank you, I'm fine."
	case "time":
		dt := time.Now()
		return dt.Format("02-01-2006 15:04:05")
	default:
		return "Unknown request"
	}
}

func init() {
	rootCmd.AddCommand(sbotCmd)
}
