package cmd

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"

	"github.com/hirosassa/zerodriver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

var (
	TeleToken   = os.Getenv("TELE_TOKEN")
	MetricsHost = os.Getenv("METRICS_HOST")
)

func initMetrics(ctx context.Context) {

	// Create a new OTLP Metric gRPC exporter with the specified endpoint and options
	exporter, _ := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(MetricsHost),
		otlpmetricgrpc.WithInsecure(),
	)

	// Define the resource with attributes that are common to all metrics.
	// labels/tags/resources that are common to all metrics.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(fmt.Sprintf("sbot_%s", appVersion)),
	)

	// Create a new MeterProvider with the specified resource and reader
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			// collects and exports metric data every 5 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(5*time.Second)),
		),
	)

	// Set the global MeterProvider to the newly created MeterProvider
	otel.SetMeterProvider(mp)

}

func pmetrics(ctx context.Context, payload string) {
	// Get the global MeterProvider and create a new Meter with the name "kbot_light_signal_counter"
	meter := otel.GetMeterProvider().Meter("sbot_light_signal_counter")

	// Get or create an Int64Counter instrument with the name "kbot_light_signal_<payload>"
	counter, _ := meter.Int64Counter(fmt.Sprintf("sbot_light_signal_%s", payload))

	// Add a value of 1 to the Int64Counter
	counter.Add(ctx, 1)
}

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
		logger := zerodriver.NewProductionLogger()

		sbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			logger.Fatal().Str("Error", err.Error()).Msg("Please check TELE_TOKEN")
			return
		} else {
			logger.Info().Str("Version", appVersion).Msg("sbot started")
		}

		sbot.Handle("/start", func(c telebot.Context) error {
			logger.Info().Str("Payload", c.Text()).Msg(c.Message().Payload)

			payload := c.Message().Payload
			pmetrics(context.Background(), payload)

			menu := &telebot.ReplyMarkup{
				ReplyKeyboard: [][]telebot.ReplyButton{
					{{Text: "hello"}, {Text: "how are you"}},
					{{Text: "time"}, {Text: "password"}, {Text: "version"}},
				},
			}
			return c.Send("Hello, I'm sbot", menu)
		})

		sbot.Handle(telebot.OnText, func(c telebot.Context) error {
			logger.Error().Str("Payload", c.Text()).Msg(c.Message().Payload)

			payload := c.Message().Payload
			pmetrics(context.Background(), payload)

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

func init() {
	ctx := context.Background()
	initMetrics(ctx)
	rootCmd.AddCommand(sbotCmd)
}
