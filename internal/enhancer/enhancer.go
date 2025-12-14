package enhancer

import (
	"context"
	"fmt"

	"promptgo/internal/ai"
)

type Input struct {
	Task       string
	Details    string
	SecretWord string
}

type Output struct {
	EnhancedPrompt string
	Tip            string
}

// QuestionsOutput represents the result of task analysis
type QuestionsOutput struct {
	TaskType  ai.TaskType
	Questions []string
}

type Enhancer struct {
	aiClient *ai.Client
}

// NewEnhancer creates a new enhancer with AI capabilities
func NewEnhancer(apiKey string, model string) *Enhancer {
	return &Enhancer{
		aiClient: ai.NewClient(apiKey, model),
	}
}

// GetQuestions analyzes the task and returns context questions (Step 1)
func (e *Enhancer) GetQuestions(ctx context.Context, task, details string) (*QuestionsOutput, error) {
	result, err := e.aiClient.AnalyzeTask(ctx, task, details)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze task: %w", err)
	}

	return &QuestionsOutput{
		TaskType:  result.TaskType,
		Questions: result.Questions,
	}, nil
}

// GeneratePrompt generates the final enhanced prompt with user answers (Step 2)
func (e *Enhancer) GeneratePrompt(ctx context.Context, input Input, taskType ai.TaskType, qa map[string]string) (*Output, error) {
	prompt, err := e.aiClient.GeneratePrompt(ctx, ai.PromptRequest{
		Task:       input.Task,
		Details:    input.Details,
		TaskType:   taskType,
		QA:         qa,
		SecretWord: input.SecretWord,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate prompt: %w", err)
	}

	tip := "This AI-generated prompt is tailored to your specific task and context. It will guide you through understanding, designing, and implementing your solution."

	return &Output{
		EnhancedPrompt: prompt,
		Tip:            tip,
	}, nil
}

// Enhance is a simple fallback function for backward compatibility
// This will be removed once TUI is updated to use the new AI flow
func Enhance(input Input) Output {
	// Simple fallback that just formats the task
	tip := "AI integration in progress. This is a temporary fallback."

	return Output{
		EnhancedPrompt: fmt.Sprintf("Task: %s\n\nDetails: %s\n\nSecret Word: %s\n\n(AI enhancement coming soon)", input.Task, input.Details, input.SecretWord),
		Tip:            tip,
	}
}
