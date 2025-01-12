# Hiking Buddies Crawler

This repository contains a web scraper written in Go that aims to collect participants' profile points on [Hiking buddies](https://www.hiking-buddies.com/) before and after taking part in a hike. The ultimate goal is to find out the underlying formula of the point gains based on the crawled data.

## Features
-	Event and Participant Data Collection: Automates the retrieval of event details, participant lists, and user points using ChromeDP for headless browser automation.
-	Gamified Point System: Tracks and updates user points based on event participation and route metrics.
-	RESTful API: Provides endpoints to retrieve event data, point gains, and worker task statuses using the Gin web framework.
-	Worker System: Executes background tasks for processing events and maintaining data consistency at configurable intervals.
-	SQLite Database: Manages persistent storage for users, events, routes, points, and credentials with modular repositories for efficient data access.
-	Robust Logging: Implements structured and leveled logging using Logrus for debugging and monitoring.

## Installation

1.	Clone the repository:

```
git clone <repository-url>
cd <repository-folder>
```

2.	Set up environment variables:
	-	HB_USERNAME: Your Hiking Buddies account email.
	-	HB_PASSWORD: Your Hiking Buddies account password.
3.	Build and run the application:

`go build ./<executable-name>`



## Usage
-	The application runs a REST API server accessible at http://localhost:8080 by default.
-	Use the following endpoints:
    -	/healthcheck: Check server health.
    -	/point-gains: Retrieve point gains data.
    -	/worker/status: View statuses of background workers.
    -	/worker/start: Start all workers.
    -	/worker/stop: Stop all workers.
