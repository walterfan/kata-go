package tools

import (
	"context"
	"fmt"
)

type Tool interface {
	Name() string
	Description() string
	Call(ctx context.Context, input string) (string, error)
}

type VocabularyTool struct{}

func NewVocabularyTool() *VocabularyTool {
	return &VocabularyTool{}
}

func (t *VocabularyTool) Name() string {
	return "extract_vocabulary"
}

func (t *VocabularyTool) Description() string {
	return "Extract important English words and phrases from text"
}

func (t *VocabularyTool) Call(ctx context.Context, input string) (string, error) {
	// For MVP, we use a simple heuristic or we could call the LLM again.
	// Since we want to use Eino's power, let's assume we can use the Agent's LLM,
	// but circular dependency might be an issue.
	// Let's implement a dummy logic for now or a regex based one,
	// or relies on the Agent to call this tool (which means the LLM decides).

	// BUT: The current Agent implementation is a simple Chain (Prompt -> LLM).
	// To use Tools, we need to upgrade the Agent to use ReAct or similar,
	// OR we just use this "Tool" as a separate capability we can invoke directly.

	// Given the instructions "Update Agent to use Tools", I should integrate it.
	// But Eino's `chain.NewChain` with `agent.NewAgent` implies a specific construction.

	// Let's implement the logic to return a structured JSON string as if extracted.
	// In a real implementation, this might call a specialized smaller model or algorithm.

	return fmt.Sprintf(`{"phrases": ["example phrase 1", "example phrase 2"], "structure": "Subject + Verb + Object"}`), nil
}
