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

func getKey() string {
	filePath := "openweatherapi.key"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "Error opening file"
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
	return apiKey
}

// func getCityLatAndLong(city string) (string, string) {
// 	fmt.Println("getCityLatAndLong: entering")
// 	apiKey := getKey()
// 	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)
// 	fmt.Println("url:", url)
// 	data := getRequest(url)
// 	fmt.Printf("getCityLatAndLong: returning: %s\n", data)
// 	fmt.Printf("getCityLatAndLong: exiting")

// func getWeatherForecast(city string, unit string) (WeatherResponse, error) {
// 	fmt.Println("getWeatherWords: entering")
// 	apiKey := getKey()
// 	url := fmt.Sprintf();
// 	fmt.Println("url:", url)
// 	weatherData = getRequest(url)
// 	fmt.Printf("getWeatherWords: returning: %s\n", weatherData.Weather[0].Main)
// 	return weatherData, nil
// }

func getRequest(url string) (WeatherResponse, error) {
	fmt.Println("getRequest: entering")
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
	fmt.Printf("getRequest: exit, returning: %s\n", weatherData.Weather[0].Main)
	return weatherData, nil
}

func getWeatherWords(city string, unit string) (string, error) {
	fmt.Println("getWeatherWords: entering")
	apiKey := getKey()
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s", city, apiKey, unit)
	fmt.Println("url:", url)
	weatherData, err := getRequest(url)
	jsonData, err := json.Marshal(weatherData)
	if err != nil {
		fmt.Printf("Error marshalling weather data: %s", err)
	}
	fmt.Printf("getWeatherWords: returning: %s\n", jsonData)
	//fmt.Printf("getWeatherWords: returning: %s\n", weatherData.Weather[0].Main)
	return weatherData.Weather[0].Main, err
}

func generateAsciiWeather(weather []byte) []byte {
	fmt.Println("Converting to ASCII art")
	asciiArt := weather
	fmt.Printf("Returning: %s\n", asciiArt)
	return asciiArt
}

func main() {
	fmt.Println("Starting app!")
	router := gin.Default()

	router.GET("/weather/:city", func(c *gin.Context) {
		// city parameter
		city := c.Param("city")
		fmt.Printf("City: %s\n", city)

		// unit query parameter
		unit := c.Query("unit")
		if unit == "" {
			unit = "imperial" // metric, imperial, standard
		}
		fmt.Printf("Unit: %s\n", unit)

		weatherData, err := getWeatherWords(city, unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		asciiWeather := generateAsciiWeather(weatherData)
		weatherString := string(asciiWeather)
		c.String(http.StatusOK, weatherString)
	})

	fmt.Println("Running service on port 8080")
	router.Run(":8080")
}
