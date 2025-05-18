package prompt

type Prompt struct {
	Name         string         `json:"name"`
	Description  string         `json:"desc"`
	SystemPrompt string         `json:"systemPrompt"`
	UserPrompt   string         `json:"userPrompt"`
	Tags         string         `json:"tags"`
}
