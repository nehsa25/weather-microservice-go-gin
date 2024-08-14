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

func getWeather(city string, unit string) (WeatherResponse, error) {
	apiKey := getKey()
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=%s", city, apiKey, unit)
	fmt.Println("url:", url)
	return getRequest(url)
}

func getWeatherJson(city string, unit string) (string, error) {
	fmt.Println("getWeatherJson: entering")
	weatherData, err := getWeather(city, unit)
	jsonData, err := json.Marshal(weatherData)
	if err != nil {
		fmt.Printf("Error marshalling weather data: %s", err)
	}
	serializedWeatherData := string(jsonData)
	fmt.Printf("getWeatherJson: returning: %s\n", serializedWeatherData)
	return serializedWeatherData, err
}

func getWeatherTemp(city string, unit string) (string, error) {
	fmt.Println("getWeatherTemp: entering")
	weatherData, err := getWeather(city, unit)
	tempString := fmt.Sprintf("%.2f", weatherData.Main.Temp)
	fmt.Printf("getWeatherTemp: returning: %s\n", tempString)
	return tempString, err
}

func getWeatherWords(city string, unit string) (string, error) {
	fmt.Println("getWeatherWords: entering")
	weatherData, err := getWeather(city, unit)
	fmt.Printf("getWeatherWords: returning: %s\n", weatherData.Weather[0].Main)
	return weatherData.Weather[0].Main, err
}

// func generateAsciiWeather(weather string) string {
// 	fmt.Println("Converting to ASCII art")
// 	art := ""
// 	switch weather {
// 	case "cloudy":
// 		art = "☁️"
// 	case "rain":
// 		art = "🌧️"
// 	case "snow":
// 		art = "❄️"
// 	case "clear":
// 		art = "☀️"
// 	case "sunny":
// 		art = "☀️"
// 	case "thunderstorm":
// 		art = "⛈️"
// 	case "mist":
// 		art = "🌫️"
// 	default:
// 		art = ""
// 	}
// 	var em = gin.H{"emoji": art}
// 	fmt.Printf("Returning: %s\n", art)
// 	return em
// }

func main() {
	fmt.Println("Starting app!")
	router := gin.Default()

	router.GET("/weather_description/:city", func(c *gin.Context) {
		city := c.Param("city")
		fmt.Printf("City: %s\n", city)

		// unit query parameter
		unit := c.Query("units")
		if unit == "" {
			unit = "imperial" // metric, imperial
		}
		fmt.Printf("Unit: %s\n", unit)

		weatherData, err := getWeatherWords(city, unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.String(http.StatusOK, string(weatherData))
	})

	router.GET("/weather_temp/:city", func(c *gin.Context) {
		city := c.Param("city")
		fmt.Printf("City: %s\n", city)

		// unit query parameter
		unit := c.Query("units")
		if unit == "" {
			unit = "imperial" // metric, imperial
		}
		fmt.Printf("Unit: %s\n", unit)

		weatherData, err := getWeatherTemp(city, unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.String(http.StatusOK, string(weatherData))
	})

	// router.GET("/weather_ascii/:city", func(c *gin.Context) {
	// 	city := c.Param("city")
	// 	fmt.Printf("City: %s\n", city)

	// 	// unit query parameter
	// 	unit := c.Query("units")
	// 	if unit == "" {
	// 		unit = "imperial" // metric, imperial, standard
	// 	}
	// 	fmt.Printf("Unit: %s\n", unit)

	// 	weatherData, err := getWeatherWords(city, unit)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	asciiWeather := generateAsciiWeather(weatherData)
	// 	c.String(http.StatusOK, asciiWeather)
	// })

	router.GET("/weather_all/:city", func(c *gin.Context) {
		city := c.Param("city")
		fmt.Printf("City: %s\n", city)

		// unit query parameter
		unit := c.Query("units")
		if unit == "" {
			unit = "imperial" // metric, imperial
		}
		fmt.Printf("Unit: %s\n", unit)

		weatherData, err := getWeatherJson(city, unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.String(http.StatusOK, string(weatherData))
	})

	fmt.Println("Running service on port 8080")
	router.Run(":8080")
}
