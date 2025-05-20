package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/viper"
)

var (
	port              string
	cachedPromptMap   map[string]interface{}
	cachedPromptList  []string // new field to preserve order
	configInitialized = false
)

func getPromptConfigByName(name string) (*PromptRequest, error) {
	item, ok := cachedPromptMap[name]
	if !ok {
		return nil, fmt.Errorf("prompt not found: %s", name)
	}

	itemMap, ok := item.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid prompt format for: %s", name)
	}

	return &PromptRequest{
		Name:         name,
		SystemPrompt: itemMap["system_prompt"].(string),
		UserPrompt:   itemMap["user_prompt"].(string),
	}, nil
}

// InitConfig initializes viper configuration only once
func InitConfig() error {
	if configInitialized {
		return nil
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	cachedPromptMap = viper.GetStringMap("prompts")
	// Preserve original order from YAML (keys are ordered in interface{})
	promptNode := viper.Get("prompts")
	if promptNode == nil {
		panic("prompts configuration is missing")
	}

	m, ok := promptNode.(map[string]interface{})
	if !ok {
		// Debugging: Print the actual type and value for diagnosis
		fmt.Fprintf(os.Stderr, "Unexpected type for 'prompts': %T, value: %+v\n", promptNode, promptNode)
		panic("prompts must be a map[string]interface{} in config")
	}

	// Proceed safely
	cachedPromptList = make([]string, 0, len(m))
	for k := range m {
		keyStr := k // Already a string, no assertion needed
		cachedPromptList = append(cachedPromptList, keyStr)
	}
	// Optional: sort by number prefix if needed
	sort.SliceStable(cachedPromptList, func(i, j int) bool {
		return cachedPromptList[i] < cachedPromptList[j]
	})
	configInitialized = true
	return nil
}
