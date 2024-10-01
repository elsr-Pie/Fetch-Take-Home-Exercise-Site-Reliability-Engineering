package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Endpoint struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
	Body    string            `yaml:"body,omitempty"`
}

type DomainStatus struct {
	TotalChecks int
	UpChecks    int
}

func loadConfig(filePath string) ([]Endpoint, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var endpoints []Endpoint
	err = yaml.Unmarshal(data, &endpoints)
	return endpoints, err
}

func checkEndpoint(endpoint Endpoint) (bool, time.Duration) {
	client := &http.Client{Timeout: 2 * time.Second}

	reqMethod := endpoint.Method
	if reqMethod == "" {
		reqMethod = "GET"
	}

	req, err := http.NewRequest(reqMethod, endpoint.URL, strings.NewReader(endpoint.Body))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return false, 0
	}

	for key, value := range endpoint.Headers {
		req.Header.Add(key, value)
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 || duration > 500*time.Millisecond {
		return false, duration
	}
	return true, duration
}

func logAvailability(statusMap map[string]DomainStatus) {
	for domain, status := range statusMap {
		availability := float64(status.UpChecks) / float64(status.TotalChecks) * 100
		fmt.Printf("%s has %.0f%% availability\n", domain, availability)
	}
}

func extractDomain(url string) string {
	urlParts := strings.Split(url, "//")
	domain := strings.Split(urlParts[1], "/")[0]
	return domain
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <config_file_path>", os.Args[0])
	}
	filePath := os.Args[1]
	endpoints, err := loadConfig(filePath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	statusMap := make(map[string]DomainStatus)

	for {
		for _, endpoint := range endpoints {
			isUp, _ := checkEndpoint(endpoint)
			domain := extractDomain(endpoint.URL)
			status := statusMap[domain]
			status.TotalChecks++
			if isUp {
				status.UpChecks++
			}
			statusMap[domain] = status
		}
		logAvailability(statusMap)
		time.Sleep(15 * time.Second)
	}
}
