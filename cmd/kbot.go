/*
Copyright © 2023 NAME HERE lapin@ucu.edu.ua
*/

package cmd


import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	telebot "gopkg.in/telebot.v3"
)

// ...

type WeatherResponse struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
}

type Weather struct {
	Description string `json:"description"`
}

type Main struct {
	Temperature float64 `json:"temp"`
	Pressure    float64 `json:"pressure"`
	Humidity    float64 `json:"humidity"`
}

// ...

func getWeather(m telebot.Context) error {
	msg := m.Message()
	cityPrompt := telebot.NewTextRequest("Введіть назву міста в Україні, для якого ви хочете дізнатися погоду: 😊🌤️")
	cityResp := m.Send(msg.Sender(), cityPrompt)

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

	weatherData, err := ioutil.ReadAll(weatherResp.Body)
	if err != nil {
		return err
	}

	weather, err := parseWeatherData(weatherData)
	if err != nil {
		return err
	}

	weatherDescription := ""
	if len(weather.Weather) > 0 {
		weatherDescription = weather.Weather[0].Description
	}

	responseText := fmt.Sprintf("Погода для міста %s:\nТемпература: %.1f°C\nТиск: %.1f гПа\nВологість: %.1f%%\nОпис: %s",
		cityName, weather.Main.Temperature, weather.Main.Pressure, weather.Main.Humidity, weatherDescription)

	return m.Send(msg.Sender(), responseText)
}

func parseWeatherData(weatherData []byte) (*WeatherResponse, error) {
	var response WeatherResponse
	err := json.Unmarshal(weatherData, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
