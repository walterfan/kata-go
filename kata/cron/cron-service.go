package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

// Define a Task structure with the 'function' and 'command' fields
type Task struct {
	Name     string `yaml:"name"`
	Schedule string `yaml:"schedule"`
	Function string `yaml:"function"`
	Command  string `yaml:"command"`
}

type Config struct {
	Tasks []Task `yaml:"tasks"`
}

// Load configuration from the YAML file
func loadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Define the function to check the URL
func check(url string) {
	// Perform a basic HTTP GET request to check if the URL is reachable
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error checking URL %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Checked URL %s with status code: %d", url, resp.StatusCode)
}

// Function to execute a command in the shell
func executeCommand(command string) {
	// Execute the command using exec.Command
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing command '%s': %v", command, err)
	}
	log.Printf("Output of command '%s': %s", command, string(output))
}

// Function to map the function name to the actual Go function
func executeFunction(functionName, param string) {
	// Remove the "()" from the function name
	functionName = strings.TrimSuffix(functionName, "()")

	log.Printf("Executing function '%s' with parameter '%s'", functionName, param)
	// Map the function name to the corresponding function
	if functionName == "check" {
		check(param)
	} else if functionName == "executeCommand" {
		executeCommand(param)
	} else {
		log.Printf("No function found for: %s", functionName)
	}
}

// Main function that reads the config and sets up the cron jobs
func main() {
	// Load configuration from the YAML file
	config, err := loadConfig("cron-config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the cron job scheduler
	c := cron.New(cron.WithSeconds()) // cron jobs with seconds precision

	// Add tasks from the YAML file to the cron scheduler
	for _, task := range config.Tasks {
		// If the task has a function, we execute it
		if task.Function != "" {
			// Parse the function and its parameters from the string
			functionDesc := task.Function
			functionDesc = strings.TrimSpace(functionDesc)

			// Extract the URL or parameter from the function string
			functionName := functionDesc
			param := ""
			if strings.Contains(functionDesc, "(") && strings.Contains(functionDesc, ")") {
				startIdx := strings.Index(functionDesc, "(") + 1
				endIdx := strings.Index(functionDesc, ")")
				if startIdx < endIdx {
					param = functionDesc[startIdx:endIdx]
					functionName = strings.TrimSpace(functionDesc[:startIdx-1])
				}
			}

			// Add the task to the cron scheduler
			_, err := c.AddFunc(task.Schedule, func(functionName, param string) func() {
				return func() {
					executeFunction(functionName, param)
				}
			}(functionName, param))

			if err != nil {
				log.Printf("Failed to add task %s: %v", task.Name, err)
			} else {
				log.Printf("Scheduled function task %s with schedule %s", task.Name, task.Schedule)
			}
		}

		// If the task has a command, we execute it
		if task.Command != "" {
			// Add the task to the cron scheduler
			_, err := c.AddFunc(task.Schedule, func(command string) func() {
				return func() {
					executeCommand(command)
				}
			}(task.Command))

			if err != nil {
				log.Printf("Failed to add task %s: %v", task.Name, err)
			} else {
				log.Printf("Scheduled command task %s with schedule %s", task.Name, task.Schedule)
			}
		}
	}

	// Start the cron scheduler
	c.Start()

	// Run indefinitely
	select {}
}
