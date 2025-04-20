package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get Redis host and password from environment variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Initialize the Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0, // Use default DB
	})
}

// Define a Task structure with the 'function' and 'command' fields
type Task struct {
	Name     string `yaml:"name"`
	Schedule string `yaml:"schedule"`
	Function string `yaml:"function"`
	Command  string `yaml:"command"`
	Deadline string `yaml:"deadline"`
}

type Config struct {
	Tasks []Task `yaml:"tasks"`
}

// TaskManager struct to encapsulate task-related functions
type TaskManager struct {
	config *Config
}

// Load configuration from the YAML file
func (tm *TaskManager) loadConfig(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return err
	}

	tm.config = &config
	return nil
}

// Define the function to check the URL
func (tm *TaskManager) check(url string) {
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
func (tm *TaskManager) executeCommand(command string) {
	// Execute the command using exec.Command
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing command '%s': %v", command, err)
	}
	log.Printf("Output of command '%s': %s", command, string(output))
}

// Function to map the function name to the actual Go function
func (tm *TaskManager) executeFunction(functionName, param string) {
	// Remove the "()" from the function name
	functionName = strings.TrimSuffix(functionName, "()")

	log.Printf("Executing function '%s' with parameter '%s'", functionName, param)
	// Map the function name to the corresponding function
	if functionName == "check" {
		tm.check(param)
	} else if functionName == "executeCommand" {
		tm.executeCommand(param)
	} else {
		log.Printf("No function found for: %s", functionName)
	}
}

// Read unfinished tasks from Redis
func (tm *TaskManager) readUnfinishedTasks() ([]string, error) {
	keys, err := rdb.Keys(ctx, "task:*").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// Publish task expiry event to Redis
func (tm *TaskManager) publishTaskExpiryEvent(taskID string) error {
	return rdb.Publish(ctx, "task_expiry_channel", taskID).Err()
}

// Check task expiry and publish events
func (tm *TaskManager) checkTaskExpiry() {
	tasks, err := tm.readUnfinishedTasks()
	if err != nil {
		log.Printf("Failed to read unfinished tasks: %v", err)
		return
	}

	for _, taskKey := range tasks {
		// Assume task expiration time is stored in Redis with key "task:<taskID>:expiry"
		expiryTime, err := rdb.Get(ctx, taskKey).Int64()
		if err != nil {
			log.Printf("Failed to get expiry time for task %s: %v", taskKey, err)
			continue
		}

		log.Printf("task %s expiry time is %d vs. %d", taskKey, expiryTime, time.Now().Unix())

		if time.Now().Unix() >= expiryTime {
			taskID := strings.TrimPrefix(taskKey, "task:")
			if err := tm.publishTaskExpiryEvent(taskID); err != nil {
				log.Printf("Failed to publish expiry event for task %s: %v", taskID, err)
			}
		}
	}
}

// Check tasks and set expiry times
func (tm *TaskManager) CheckTasks() {
	// Load configuration from the YAML file
	if err := tm.loadConfig("cron-config.yml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the cron job scheduler
	c := cron.New(cron.WithSeconds()) // cron jobs with seconds precision

	// Add tasks from the YAML file to the cron scheduler
	for _, task := range tm.config.Tasks {
		// If the task has a function, we execute it
		if task.Function != "" && task.Schedule != "" {
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
					tm.executeFunction(functionName, param)
				}
			}(functionName, param))

			if err != nil {
				log.Printf("Failed to add task %s: %v", task.Name, err)
			} else {
				log.Printf("Scheduled function task %s with schedule %s", task.Name, task.Schedule)
			}
		}

		// If the task has a command, we execute it
		if task.Command != "" && task.Schedule != "" {
			// Add the task to the cron scheduler
			_, err := c.AddFunc(task.Schedule, func(command string) func() {
				return func() {
					tm.executeCommand(command)
				}
			}(task.Command))

			if err != nil {
				log.Printf("Failed to add task %s: %v", task.Name, err)
			} else {
				log.Printf("Scheduled command task %s with schedule %s", task.Name, task.Schedule)
			}
		}

		if task.Deadline != "" {
			// Parse the deadline string to a time.Time object
			deadlineTime, err := time.Parse(time.RFC3339, task.Deadline)
			if err != nil {
				log.Printf("Failed to parse deadline for task %s: %v", task.Name, err)
				continue
			}

			// Convert the deadline time to a Unix timestamp
			expiryTime := deadlineTime.Unix()

			// Set the Redis key with the expiry time
			key := fmt.Sprintf("task:%s:expiry", task.Name)
			key = strings.ReplaceAll(key, " ", "_")

			err = rdb.Set(ctx, key, expiryTime, 0).Err()
			if err != nil {
				log.Printf("Failed to set expiry time for task %s: %v", task.Name, err)
				continue
			}

			log.Printf("Set expiry time for task %s: %d", task.Name, expiryTime)
		}
	}

	// Add the task expiry check cron job
	_, err := c.AddFunc("@every 1m", tm.checkTaskExpiry)
	if err != nil {
		log.Fatalf("Failed to add task expiry check cron job: %v", err)
	} else {
		log.Println("Scheduled task expiry check cron job")
	}

	// Start the cron scheduler
	c.Start()

	// Run indefinitely
	select {}
}

func main() {
	// Initialize the TaskManager
	tm := &TaskManager{}

	// Check tasks and set expiry times
	tm.CheckTasks()
}
