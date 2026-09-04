package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/IvanPeterschmitt/weather-cli/internal/api"
	"github.com/IvanPeterschmitt/weather-cli/internal/cli"
)



func main() {
	args := os.Args[1:]
	
	if len(args) == 0{
		fmt.Fprintln(os.Stderr, "No city mentioned. Please, enter a city as an argument in the command line.")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Too many arguments. Please, only mention one city.\nIf the city name is composed of several words, write it between quotes.")
		os.Exit(1)
	}

	cityName := args[0]

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Note: no .env file found, system variable reading")
	}

	client := api.NewClient(os.Getenv("API_TOKEN"))

	citiesInformation, err := api.GetCityCode(cityName, client)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while city code retrieval: %v\n", err)
		os.Exit(1)
	}

	var cityIndex int

	if len(citiesInformation) > 1 {
		cityIndex = cli.DisplayChoice(citiesInformation)
	}

	weatherData, err := api.GetWeatherData(&citiesInformation[cityIndex], client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while weather data retrieval: %v\n", err)
		os.Exit(1)
	}
	if weatherData == nil {
		fmt.Fprintf(os.Stderr, "Null pointer retrieved from weather data function: %v\n", err)
		os.Exit(1)
	}

	cli.DisplayResult(citiesInformation[cityIndex].Name, weatherData)
}