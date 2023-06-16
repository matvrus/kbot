/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
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
	// Telegram bot token
	TeleToken = os.Getenv("TELE_TOKEN")
	// MetricsHost exporter host:port
	MetricsHost = os.Getenv("METRICS_HOST")
)

// Initialize OpenTelemetry for metrics
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
	resourceAttributes := []attribute.KeyValue{
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(fmt.Sprintf("kbot_%s", appVersion)),
	}

	resource := resource.NewWithAttributes(resourceAttributes...)

	// Create a new MeterProvider with the specified resource and exporter
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithBatcher(exporter),
	)

	// Register the provider globally
	otel.SetMeterProvider(provider)
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
		ctx := context.Background()
		initMetrics(ctx)

		fmt.Println("kbot %s started", appVersion)
		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
			return
		}

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			log.Print(m.Message().Payload, m.Text())
			payload := m.Message().Payload

			switch payload {
			case "hello":
				err = m.Send(fmt.Sprintf("Hello, %s! I'm Kbot %s!", m.Sender().FirstName, appVersion))
			}

			return err
		})

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
}
