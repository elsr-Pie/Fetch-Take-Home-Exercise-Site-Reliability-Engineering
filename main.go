package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
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
	ErrorBudget float64
}

var slo float64 = 97.0  // Setting SLO to 97%
var errorBudget float64 // Automatically calculated as 100 - SLO

func init() {
	// Automate error budget calculation
	errorBudget = 100 - slo
}

func main() {
	// Load configuration file from argument
	if len(os.Args) < 2 {
		log.Fatal("Please provide the YAML configuration file path")
	}
	configFile := os.Args[1]

	// Read and parse the YAML file
	fileContent, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	var endpoints []Endpoint
	err = yaml.Unmarshal(fileContent, &endpoints)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	// Initialize domain status map
	domainStatusMap := make(map[string]*DomainStatus)

	for {
		var wg sync.WaitGroup

		for _, endpoint := range endpoints {
			wg.Add(1)
			// Run each endpoint check in parallel
			go func(ep Endpoint) {
				defer wg.Done()
				status := checkEndpoint(ep)
				updateDomainStatus(domainStatusMap, ep.URL, status)
			}(endpoint)
		}

		// Wait for all Goroutines to complete
		wg.Wait()

		// Log the availability after each round of checks
		logAvailability(domainStatusMap)

		// Wait for 15 seconds before the next check
		time.Sleep(15 * time.Second)
	}
}

// checkEndpoint sends the HTTP request to the endpoint and checks the status
func checkEndpoint(endpoint Endpoint) bool {
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	// Set up the HTTP request
	req, err := http.NewRequest(endpoint.Method, endpoint.URL, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v", endpoint.URL, err)
		return false
	}

	// Add headers if provided
	for key, value := range endpoint.Headers {
		req.Header.Add(key, value)
	}

	// Perform the HTTP request
	startTime := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		log.Printf("Request failed for %s: %v", endpoint.URL, err)
		return false
	}
	defer resp.Body.Close()

	// Check if the response is UP (2xx status and < 500ms response time)
	if resp.StatusCode >= 200 && resp.StatusCode < 299 && duration < 500*time.Millisecond {
		return true
	}

	return false
}

// updateDomainStatus updates the domain's availability and error budget
func updateDomainStatus(statusMap map[string]*DomainStatus, url string, up bool) {
	domain := extractDomain(url)

	if _, exists := statusMap[domain]; !exists {
		statusMap[domain] = &DomainStatus{}
	}

	statusMap[domain].TotalChecks++

	if up {
		statusMap[domain].UpChecks++
	} else {
		// Error budget is consumed when an endpoint check fails
		statusMap[domain].ErrorBudget += 100.0 / float64(statusMap[domain].TotalChecks)
	}
}

// logAvailability logs the current availability and error budget for each domain
func logAvailability(statusMap map[string]*DomainStatus) {
	for domain, status := range statusMap {
		availability := 100.0 * float64(status.UpChecks) / float64(status.TotalChecks)

		fmt.Printf("%s has %d%% availability\n", domain, int(availability))

		// Check if SLO is being violated
		if availability < slo {
			fmt.Printf("Warning: %s is below the SLO of %.1f%%\n", domain, slo)
		}

		// Error budget usage
		fmt.Printf("%s has consumed %.2f%% of its error budget\n", domain, status.ErrorBudget)
		if status.ErrorBudget >= errorBudget {
			fmt.Printf("Alert: %s has exhausted its error budget!\n", domain)
		}
	}
}

// extractDomain gets the domain name from a URL
func extractDomain(url string) string {
	// Assuming a basic split by "//" to get the domain name
	// For a more complex solution, consider using the "net/url" package
	split := len("https://")
	return url[split:]
}
