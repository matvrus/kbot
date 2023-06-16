/*
Copyright © 2023 NAME HERE lapin@ucu.edu.ua
*/

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	telebot "gopkg.in/tucnak/telebot.v3"
)

var (
	// Telegram bot token
	TeleToken = os.Getenv("TELE_TOKEN")
	// OpenWeatherMap API key
	WeatherAPIKey = os.Getenv("WEATHER_API_KEY")
	// MetricsHost exporter host:port
	MetricsHost = os.Getenv("METRICS_HOST")
)

// Initialize OpenTelemetry
func initMetrics(ctx context.Context) {

	// Create a new OTLP Metric gRPC exporter with the specified endpoint and options
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(MetricsHost),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Define the resource with attributes that are common to all metrics.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(fmt.Sprintf("kbot_%s", appVersion)),
	)

	// Create a new MeterProvider with the specified resource and exporter
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithSyncer(exporter),
	)
	otel.SetMeterProvider(mp)
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
		logger := log.New(os.Stdout, "", log.LstdFlags)

		ctx := context.Background()
		initMetrics(ctx)

		logger.Printf("kbot %s started\n", appVersion)

		bot, err := telebot.NewBot(telebot.Settings{
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})
		if err != nil {
			logger.Fatalf("Please check TELE_TOKEN env variable. %s", err)
		}

		bot.Handle(telebot.OnText, func(m *telebot.Message) error {
			logger.Println(m.Payload, m.Text)

			switch m.Text {
			case "hello":
				err := m.Reply("world")
				if err != nil {
					logger.Println("Error:", err)
				}
			case "/start":
				helpText := "Доступні команди:\n" +
					"/start - Початок роботи\n" +
					"/help - Довідка\n" +
					"/echo - Ехо-відповідь\n" +
					"/time - Поточний час\n" +
					"/weather - Погода в Україні"
				err = m.Reply(helpText)
			case "/help":
				helpText := "Доступні команди:\n" +
					"/help - Довідка\n" +
					"/echo - Ехо-відповідь\n" +
					"/time - Поточний час\n" +
					"/weather - Погода в Україні"
				err = m.Reply(helpText)
			case "/echo":
				text := m.Text
				err = m.Reply(text)
			case "/time":
				currentTime := time.Now().Format("2006-01-02 15:04:05")
				err = m.Reply(fmt.Sprintf("Поточний час: %s ⌚", currentTime))
			case "/weather":
				err = getWeather(m)
			default:
				err = m.Reply("Не розумію вашої команди. Введіть /help для довідки. 😕")
			}

			return err
		})

		bot.Start()
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
}

func getWeather(m telebot.Context) error {
	msg := m.Message()

	cityPrompt := telebot.NewTextRequest("Введіть назву міста в Україні, для якого ви хочете дізнатися погоду: 😊🌤️")
	cityResp := m.Send(msg.Sender, cityPrompt)

	cityName := ""

	for cityResp.Next() {
		cityName = cityResp.Text()
		break
	}

	if cityName == "" {
		return nil
	}

	weatherURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s,ua&appid=%s&units=metric", cityName, WeatherAPIKey)
	weatherResp, err := http.Get(weatherURL)
	if err != nil {
		return err
	}
	defer weatherResp.Body.Close()

	_, err = ioutil.ReadAll(weatherResp.Body)
	if err != nil {
		return err
	}

	return m.Send(msg.Sender, "Отримано інформацію про погоду для міста "+cityName+"! 🌤️")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
