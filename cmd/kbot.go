package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/hirosassa/zerodriver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	telebot "gopkg.in/telebot.v3"
)

var (
	// TeleToken bot
	TeleToken   = os.Getenv("TELE_TOKEN")
	// MetricsHost exporter host:port
	MetricsHost = os.Getenv("METRICS_HOST")
)

// Initialize OpenTelemetry
func initMetrics(ctx context.Context) {
	// Create a new OTLP Metric gRPC exporter with the specified endpoint and options
	exporter, _ := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(MetricsHost),
		otlpmetricgrpc.WithInsecure(),
	)

	// Define the resource with attributes that are common to all metrics.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(fmt.Sprintf("kbot_%s", appVersion)),
	)

	// Create a new MeterProvider with the specified resource and reader
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(10*time.Second)),
		),
	)

	// Set the global MeterProvider to the newly created MeterProvider
	otel.SetMeterProvider(mp)
}

// HandleTelegramCommand handles different commands received from Telegram
func HandleTelegramCommand(m telebot.Context) error {
	payload := m.Message().Payload

	switch payload {
	case "hello":
		err := m.Send(fmt.Sprintf("Hello, I'm Kbot %s!", appVersion))
		return err
	case "/help":
		helpText := "Доступні команди:\n" +
			"/hello - Привітання\n" +
			"/help - Довідка\n" +
			"/echo - Ехо-відповідь\n" +
			"/time - Поточний час\n" +
			"/weather - Погода в Україні"
		err := m.Send(helpText)
		return err
	case "/echo":
		text := m.Text()
		err := m.Send(text)
		return err
	case "/time":
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		err := m.Send(fmt.Sprintf("Поточний час: %s ⌚", currentTime))
		return err
	case "/weather":
		weatherText := getWeather()
		err := m.Send(weatherText)
		return err
	default:
		err := m.Send("Не розумію вашої команди. Введіть /help для довідки. 😕")
		return err
	}
}

// getWeather видає актуальну інформацію про погоду в Україні
func getWeather() string {
	// Реалізуйте логіку отримання погоди тут
	// Поверніть актуальну інформацію про погоду у форматі string
	weatherText := "Погода в Україні: сонячно 🌞"
	return weatherText
}

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerodriver.NewProductionLogger()

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			logger.Fatal().Str("Error", err.Error()).Msg("Please check TELE_TOKEN")
			return
		} else {
			logger.Info().Str("Version", appVersion).Msg("kbot started")
		}

		kbot.Handle(telebot.OnText, HandleTelegramCommand)

		kbot.Start()
	},
}

func init() {
	ctx := context.Background()
	initMetrics(ctx)
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Initialize OpenTelemetry tracer
}
