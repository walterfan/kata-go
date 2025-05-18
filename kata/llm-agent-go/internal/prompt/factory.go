// internal/prompt/factory.go
package prompt

import (
	"github.com/goccy/go-yaml"
	"os"
)

var prompts map[string]Prompt

// LoadPrompts loads the prompts from a YAML file into the prompts map
func LoadPrompts(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var promptList []Prompt
	err = yaml.Unmarshal(data, &promptList)
	if err != nil {
		return err
	}

	prompts = make(map[string]Prompt)
	for _, p := range promptList {
		prompts[p.Name] = p
	}

	return nil
}

// GetPromptByName returns a Prompt by its name, or an error if not found
func GetPromptByName(name string) (Prompt, error) {
	prompt, exists := prompts[name]
	if !exists {
		return Prompt{}, fmt.Errorf("prompt with name '%s' not found", name)
	}
	return prompt, nil
}