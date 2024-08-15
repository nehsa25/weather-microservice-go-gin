# Weather Microservice in Go with Gin

This document outlines a basic weather microservice written in Golang using the Gin framework.

# Features:

Retrieves weather information for a specified city.
Optionally accepts a unit parameter (metric or imperial).
Exposes a GET API endpoint for integration.
Plan is to eventually have it also be able to return as ascii and emojis but not yet..

# Getting Started:

Prerequisites:

Go installed (https://go.dev/doc/install)
Basic understanding of Go and Gin

# Installation:

- Download the code (replace with your download method)
- Navigate to the project directory:
- cd weather-microservice

# Run the service:
- go run main.go

This starts the microservice on port 8080 by default (modifiable in the code).

# Usage:

The service exposes a GET endpoint at /weather/:city (replace :city with the desired city name). You can optionally specify a unit (metric or imperial) in the query parameter unit:

http://localhost:8080/weather_all?city=Vancouver?unit=metric
http://localhost:8080/weather_description?city=Vancouver?unit=metric
http://localhost:8080/weather_temp?city=Vancouver?unit=imperial

This will return weather information for Vancouver in JSON format.

# Configuration:
None

# Development:
dockerfile provided

All contributions welcome to this project. Feel free to submit pull requests!

License:

This project is licensed under the MIT License (see LICENSE file for details).
