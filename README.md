# Fetch-Take-Home-Exercise-Site-Reliability-Engineering

Welcome to the comprehensive README for the HTTP Endpoint Health Check Program!

This guide will walk you through the process of setting up and running the solution, with a detailed explanation of what the Go program does.

Table of Contents

1. What You'll Need
2. Getting Started]
3. Let's Run the Program
4. Understanding the Code
5. What You'll See

What You'll Need

To get started, please ensure you have the following on your machine:

1. Go programming language**: Version 1.22 or later is recommended. You can download and install Go from the official website `https://go.dev/dl/`
2. Git: A version control system that helps you clone repositories or manage your own.


Getting Started

1. Clone the Repository:
   
   Open a terminal and run the following command to clone the repository containing the project:
   
   `git clone https://github.com/your-repo-name/health-checker.git`
   
   Alternatively, create a directory and manually place the Go code (`main.go`) and the YAML configuration file (`http-endpoints.yaml`) inside that directory.

2. Install Go:
   
   If Go is not installed, download and install it from "https://go.dev/dl/". Once installed, verify the installation by running the following command:
   
   `go version`
   
   You should see output similar to:
   
   `go version go1.22.4 linux/amd64`
   

Let's Get the Program Up and Running!

Creating Your YAML Configuration File

To start, you'll need to set up a YAML configuration file that contains the HTTP endpoints you want to monitor. For instance, create a file named `http-endpoint.yaml` like the one I have craeted or give it any name you prefer in the same directory as your Go program with content similar to what I have created in http-endpoint.yaml file.


Running or Building the Program

You have two options for running the program: executing it directly or compiling and running the compiled version.

#### Option 1: Run Directly

Open the terminal in the project directory and execute the following command:

`go run main.go endpoints.yaml`


This will run the program immediately and initiate health checks on the endpoints listed in `http-endpoint.yaml` or the yaml file you created.

Option 2: Build and Run

First, compile the program:

`go build main.go`


This will generate an executable file (e.g., `main` or `main.exe` on Windows). Next, execute the compiled program with:


`./main endpoints.yaml`


Seeing the Results

The program will print the availability percentages of each domain every 15 seconds. Example console output:

`
www.fetchrewards.com has 100% availability
fetch.com has 67% availability
fetch.com has 67% availability
www.fetchrewards.com has 100% availability
`

Stopping the Program

To terminate the program, press `CTRL+C` in the terminal running the program.

Let's Explore the Code Together!

1. The Starting Point: Main Function

The `main()` function is where the program begins. It checks if you've provided the correct argument (the path to the YAML configuration file) and exits with an error message if not. This ensures you have the necessary input before proceeding.

```go
if len(os.Args) != 2 {
   log.Fatalf("Usage: %s <config_file_path>", os.Args[0])
}
```

2. Loading the Configuration**

The `loadConfig` function reads your YAML file and converts it into a Go data structure called `[]Endpoint`. Each endpoint has fields like `name`, `url`, `method`, `headers`, and `body`.

```go
endpoints, err := loadConfig(filePath)
if err != nil {
   log.Fatalf("Failed to load config: %v", err)
}
```

3. Performing HTTP Health Checks

The `checkEndpoint()` function sends an HTTP request to each endpoint, measuring response time and checking for a 2xx status code. It considers an endpoint "UP" if it responds quickly (under 500ms) and has a valid status code.

```go
if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 || duration > 500*time.Millisecond {
   return false, duration
}
```

4. Keeping Track of Availability

The program uses a `DomainStatus` structure to track each domain's health status. It records the total checks and "UP" responses per domain.

```go
status.TotalChecks++
if isUp {
   status.UpChecks++
}
```

5. Logging Results

Every 15 seconds, the program calculates and logs the availability percentage of each domain based on checks performed so far.

```go
availability := float64(status.UpChecks) / float64(status.TotalChecks) * 100
fmt.Printf("%s has %.0f%% availability\n", domain, availability)
```

6. Continuous Monitoring

The program runs in an endless loop, checking endpoints every 15 seconds until you press `CTRL+C`.

```go
time.Sleep(15 * time.Second)
```


Anticipated Output

The program prints the availability percentage of each domain in the console every 15 seconds. An example output might be:

```
www.fetchrewards.com has 100% availability
fetch.com has 67% availability
fetch.com has 67% availability
www.fetchrewards.com has 100% availability
fetch.com has 67% availability
www.fetchrewards.com has 100% availability
```

This shows how often each domain's endpoints are "UP" during the program's runtime.


I hope this friendly overview has helped you understand how to set up, run, and appreciate the HTTP Endpoint Health Check program written in Go. Happy monitoring! 😊
