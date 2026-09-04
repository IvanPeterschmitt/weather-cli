# Weather CLI

A modular and efficient Command Line Interface (CLI) tool written in **Go**, designed to fetch weather forecasts for French municipalities using the [Météo-Concept API](https://www.meteo-concept.com/).

## Features

* **Flexible City Input**: Seamlessly handles single-word and multi-word city names without quotes (e.g., `Saint-Étienne`, `"Saint Etienne"`).
* **Interactive Disambiguation**: When a search query matches multiple cities, the CLI prompts an interactive selection menu directly in the terminal.
* **Safe Networking**: Built on Go's native `net/url` package for proper URL building and encoding, ensuring full protection against special character bugs and query injection.
* **Clean Architecture**: Strict one-way dependency flow (`main` -> `cli` -> `api`) preventing circular imports and isolating business logic, network calls, and user interactions.

---

## Architecture & Project Structure

The project follows clean architecture principles to ensure testability, maintainability, and clear separation of concerns:

```text
weather-cli/
├── cmd/
│   └── weather/
│       └── main.go          # Entry point & Controller (Orchestrator)
├── internal/
│   ├── api/                 # Data fetching & Network layer
│   │   ├── client.go        # Generic HTTP client & query parameter handling
│   │   ├── geocoding.go     # City lookup & INSEE code resolution
│   │   └── weather.go       # Weather API structs & endpoints
│   └── cli/                 # Terminal UI & Interactive user input
│       └── cli.go           # Terminal prompts, selection logic & output formatting
├── .env                     # Local environment configuration (API Token)
├── go.mod
└── README.md