package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

var (
	TeleToken = os.Getenv("TELE_TOKEN")
)

func init() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rootCmd.AddCommand(sbotCmd)
}

var sbotCmd = &cobra.Command{
	Use:     "sbot",
	Aliases: []string{"start"},
	Short:   "Start telegram sbot application",
	Long: `A simple Telegram bot that can handle text messages.
	You can write something to https://t.me/yuandrk_bot and sometimes it answers :)`,
	Run: func(cmd *cobra.Command, args []string) {
		sbot, err := telebot.NewBot(telebot.Settings{
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Error creating bot, please check TELE_TOKEN env variable: %s", err)
			return
		}

		sbot.Handle("/start", func(c telebot.Context) error {
			menu := &telebot.ReplyMarkup{
				ReplyKeyboard: [][]telebot.ReplyButton{
					{{Text: "hello"}, {Text: "how are you"}},
					{{Text: "time"}, {Text: "password"}, {Text: "version"}},
				},
			}
			return c.Send("Hello, I'm sbot", menu)
		})

		sbot.Handle(telebot.OnText, func(c telebot.Context) error {
			response := handlePayload(c.Message().Text)
			err := c.Send(response) // Send returns only error
			if err != nil {
				log.Println("Error sending message:", err)
			}
			return err
		})

		log.Printf("sbot %s started\n", appVersion)
		sbot.Start()
	},
}

func handlePayload(payload string) string {
	switch strings.ToLower(payload) {
	case "hello":
		return fmt.Sprintf("Hello, I'm sbot %s", appVersion)
	case "how are you":
		return "Thank you, I'm fine."
	case "time":
		dt := time.Now()
		return dt.Format("02-01-2006 15:04:05")
	case "password":
		return fmt.Sprintf("Your random password: %s", generatePassword())
	case "version":
		return fmt.Sprintf("Current sbot version: %s", appVersion)
	default:
		return "Unknown request"
	}
}

func generatePassword() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var numbers = []rune("0123456789")
	var specialChars = []rune("!@#$%^&*()")

	password := make([]rune, 8)
	for i := range password {
		if i == 0 {
			password[i] = specialChars[rand.Intn(len(specialChars))]
		} else if rand.Intn(2) == 0 {
			password[i] = letters[rand.Intn(len(letters))]
		} else {
			password[i] = numbers[rand.Intn(len(numbers))]
		}
	}
	return string(password)
}
