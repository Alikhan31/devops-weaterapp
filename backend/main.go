package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type WeatherResponse struct {
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	City        string  `json:"city"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"windSpeed"`
}

type OpenWeatherResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

var redisClient *redis.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	} else {
		log.Println(".env file loaded successfully")
	}

	// Debug: Print all environment variables
	apiKey := getEnv("OPENWEATHER_API_KEY", "")
	if apiKey == "" {
		log.Println("WARNING: OPENWEATHER_API_KEY is not set")
	} else {
		log.Printf("OPENWEATHER_API_KEY is set (first 4 chars: %s...)", apiKey[:4])
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
		Password: "",
		DB:       0,
	})
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getWeatherFromAPI(lat, lon string) (*WeatherResponse, error) {
	apiKey := getEnv("OPENWEATHER_API_KEY", "")
	if apiKey == "" {
		return nil, fmt.Errorf("OpenWeather API key not found in environment variables")
	}

	// Print API key for debugging (first 4 characters)
	log.Printf("Using API key: %s...", apiKey[:4])

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s&units=metric", lat, lon, apiKey)
	log.Printf("Requesting URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP request error: %v", err)
		return nil, fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("API Error Response: %s", string(body))
		return nil, fmt.Errorf("weather API returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("API Response: %s", string(body))

	var weatherData OpenWeatherResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return nil, fmt.Errorf("failed to decode weather data: %v", err)
	}

	return &WeatherResponse{
		Temperature: weatherData.Main.Temp,
		Description: weatherData.Weather[0].Description,
		City:        weatherData.Name,
		Humidity:    weatherData.Main.Humidity,
		WindSpeed:   weatherData.Wind.Speed,
	}, nil
}

func getWeather(c *gin.Context) {
	lat := c.Query("lat")
	lon := c.Query("lon")

	if lat == "" || lon == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing coordinates"})
		return
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("weather:%s:%s", lat, lon)
	cachedWeather, err := redisClient.Get(c, cacheKey).Result()
	if err == nil {
		var weather WeatherResponse
		if err := json.Unmarshal([]byte(cachedWeather), &weather); err == nil {
			c.JSON(http.StatusOK, weather)
			return
		}
	}

	// If not in cache, get from API
	weather, err := getWeatherFromAPI(lat, lon)
	if err != nil {
		log.Printf("Error getting weather: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cache the result for 1 hour
	weatherJSON, _ := json.Marshal(weather)
	redisClient.Set(c, cacheKey, weatherJSON, time.Hour)

	c.JSON(http.StatusOK, weather)
}

func healthCheck(c *gin.Context) {
	// Check Redis connection
	ctx := c.Request.Context()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"error":  "Redis connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func main() {
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", healthCheck)
	r.GET("/weather", getWeather)

	log.Fatal(r.Run(":8080"))
}
