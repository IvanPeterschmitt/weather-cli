package cli
import (
	"fmt"
	"os"
	"github.com/IvanPeterschmitt/weather-cli/internal/api"
)


func DisplayChoice(choices []api.City) int {
	var cityIndex int
	var indexErr bool = true

	fmt.Fprintln(os.Stdout, "Various cities match your input:")

	for i, city := range choices {
		fmt.Fprintf(os.Stdout, "%d. %s\n", i, city.Name)
	}

	for indexErr == true {
		fmt.Fprintf(os.Stdout, "\nWhich one is the correct city? (0, %d) ", len(choices)-1)
		fmt.Scanln(&cityIndex)
		//valider la saisi
		if (int(cityIndex) < 0 || int(cityIndex) >= len(choices)) {
			fmt.Fprint(os.Stdout, "Index error: selected city index is out of range")
		} else {
			indexErr = false
		}
	}

	return cityIndex
}

func DisplayResult(cityName string, results *api.WeatherData) {
	fmt.Fprintln(os.Stdout, "\n" + cityName)
	fmt.Fprintln(os.Stdout, "────────────────────────────")
	fmt.Fprintln(os.Stdout, "Temperature  " + results.CurrentData.Measures.Temperature.Celcius + "°C")
	fmt.Fprintln(os.Stdout, "Humidity  " + results.CurrentData.Measures.Humidity.Percent + "%")
	fmt.Fprintln(os.Stdout, "Wind  " + results.CurrentData.Measures.Wind.Speed + " km/h\n")
	fmt.Fprintln(os.Stdout, "Forecast:")
	fmt.Fprintf(os.Stdout, "Today            %d°C / %d°C\n", results.ForecastData.Day[0].Data.Tmin, results.ForecastData.Day[0].Data.Tmax)
	fmt.Fprintf(os.Stdout, "Tomorrow         %d°C / %d°C\n", results.ForecastData.Day[1].Data.Tmin, results.ForecastData.Day[1].Data.Tmax)
	fmt.Fprintf(os.Stdout, "After Tomorrow   %d°C / %d°C\n", results.ForecastData.Day[2].Data.Tmin, results.ForecastData.Day[2].Data.Tmax)
}