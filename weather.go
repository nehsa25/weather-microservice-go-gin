package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
}

func getWeather(city string, unit string) (WeatherResponse, error) {
	fmt.Println("getWeather: entering")
	filePath := "openweatherapi.key"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return WeatherResponse{}, err
	}
	defer file.Close()

	// Read the API key from the file
	scanner := bufio.NewScanner(file)
	var apiKey string
	if scanner.Scan() {
		apiKey = scanner.Text()
		fmt.Println("API Key:", apiKey)
	} else {
		fmt.Println("Error reading API key:", scanner.Err())
	}
	fmt.Println("API Key2:", apiKey)
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s", city, apiKey, unit)
	fmt.Println("url:", url)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, err
	}

	var weatherData WeatherResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return WeatherResponse{}, err
	}

	fmt.Printf("getWeather: returning: %s\n", weatherData.Weather[0].Main)

	return weatherData, nil
}

func generateAsciiWeather(weather string) string {
	fmt.Println("Converting to ASCII art")
	asciiArt := weather
	fmt.Printf("Returning: %s\n", asciiArt)
	return asciiArt
}

func main() {
	fmt.Println("Starting app!")
	router := gin.Default()

	router.GET("/weather/:city/:unit", func(c *gin.Context) {
		city := c.Param("city")
		fmt.Printf("City: %s\n", city)
		unit := c.Param("unit")
		fmt.Printf("Unit: %s\n", unit)
		weatherData, err := getWeather(city, unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		asciiWeather := generateAsciiWeather(weatherData.Weather[0].Main)
		c.String(http.StatusOK, asciiWeather)
	})

	fmt.Println("Running service on port 8080")
	router.Run(":8080")
}
