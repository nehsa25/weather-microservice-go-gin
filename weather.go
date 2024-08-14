package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		Id          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		Id      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
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

// sendGetRequest makes a GET request to the given URL with the given parameters and returns the response as a WeatherResponse struct.
func sendGetRequest(urlString string, params map[string]string) (WeatherResponse, error) {
	fmt.Println("getRequest: entering")

	// Create a URL with query parameters
	u, err := url.Parse(urlString)
	if err != nil {
		// Handle error
		fmt.Println("Error parsing URL:", err)
		return WeatherResponse{}, err
	}

	q := u.Query()
	q.Set("q", params["q"])
	q.Set("appid", params["appid"])
	q.Set("units", params["units"])
	u.RawQuery = q.Encode()

	finalURL := u.String()
	fmt.Println("Final URL:", finalURL)

	// Make the GET request
	resp, err := http.Get(finalURL)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, err
	}

	// check status
	if resp.StatusCode != http.StatusOK {
		return WeatherResponse{}, fmt.Errorf("request was not successful: %s", body)
	}

	var weatherData WeatherResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return WeatherResponse{}, err
	}
	fmt.Printf("getRequest: exit, returning: %s\n", weatherData.Weather[0].Main)
	return weatherData, nil
}

// getWeatherv25 is a handler function to fetch weather data for a given city and unit using openweathermaps v2.5 API.
// It should return the weather data as a byte slice and an error if something goes wrong.
func getWeatherv25(city string, units string) (WeatherResponse, error) {
	params := map[string]string{
		"q":     city,
		"appid": getKey(),
		"units": units,
	}
	url := "https://api.openweathermap.org/data/2.5/weather"
	response, err := sendGetRequest(url, params)
	return response, err
}

func getWeatherJson(city string, unit string) (string, error) {
	fmt.Println("getWeatherJson: entering")
	weatherData, err := getWeatherv25(city, unit)
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
	weatherData, err := getWeatherv25(city, unit)
	tempString := fmt.Sprintf("%.2f", weatherData.Main.Temp)
	fmt.Printf("getWeatherTemp: returning: %s\n", tempString)
	return tempString, err
}

func getWeatherWords(city string, unit string) (string, error) {
	fmt.Println("getWeatherWords: entering")
	weatherData, err := getWeatherv25(city, unit)
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

	router.GET("/weather_description", func(c *gin.Context) {
		// get city parameter
		city := c.Query("city")
		if city == "" {
			city = "Cape Canaveral, FL" // "city" or "city, state" or "city, state, country"
		}
		fmt.Printf("city: %s\n", city)

		// units query parameter
		units := c.Query("units")
		if units == "" {
			units = "imperial" // metric, imperial
		}
		fmt.Printf("units: %s\n", units)

		// makes the GET request to the OpenWeatherMap API
		weatherData, err := getWeatherWords(city, units)
		if err != nil {
			fmt.Printf("router weather_description: Error getting weather data: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.String(http.StatusOK, string(weatherData))
	})

	router.GET("/weather_temp", func(c *gin.Context) {
		// get city parameter
		city := c.Query("city")
		if city == "" {
			city = "Cape Canaveral, FL" // "city" or "city, state" or "city, state, country"
		}
		fmt.Printf("city: %s\n", city)

		// units query parameter if provided (default is imperial)
		units := c.Query("units")
		if units == "" {
			units = "imperial" // metric, imperial
		}
		fmt.Printf("units: %s\n", units)

		// makes the GET request to the OpenWeatherMap API
		weatherData, err := getWeatherTemp(city, units)
		if err != nil {
			fmt.Printf("router weather_temp: Error getting weather data: %s", err)
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

	router.GET("/weather_all", func(c *gin.Context) {
		// get city parameter
		city := c.Query("city")
		if city == "" {
			city = "Cape Canaveral, FL, United States" // "city" or "city, state" or "city, state, country"
		}
		fmt.Printf("city: %s\n", city)

		// units query parameter
		units := c.Query("units")
		if units == "" {
			units = "imperial" // metric, imperial
		}
		fmt.Printf("units: %s\n", units)

		// makes the GET request to the OpenWeatherMap API
		weatherData, err := getWeatherJson(city, units)
		if err != nil {
			fmt.Printf("router weather_all: Error getting weather data: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.String(http.StatusOK, string(weatherData))
	})

	fmt.Println("Running service on port 8080")
	router.Run(":8080")
}
