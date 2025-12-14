package ai

import (
	"context"
	"encoding/json"
	"fmt"
)

type TaskType string

const (
	TypeFeature       TaskType = "feature"
	TypeBugFix        TaskType = "bugfix"
	TypeTesting       TaskType = "testing"
	TypeRefactoring   TaskType = "refactoring"
	TypeDocumentation TaskType = "documentation"
	TypeOther         TaskType = "other"
)

type AnalysisResult struct {
	TaskType  TaskType `json:"task_type"`
	Questions []string `json:"questions"`
}

// AnalyzeTask analyzes a task and generates context-gathering questions
func (c *Client) AnalyzeTask(ctx context.Context, task, details string) (*AnalysisResult, error) {
	systemPrompt := `You are an expert software development assistant analyzing a developer's task.

Your job:
1. Classify the task type (feature, bugfix, testing, refactoring, documentation, or other)
2. Generate 2-4 intelligent follow-up questions to gather context

Questions should be:
- Specific to the task type
- Help understand constraints, existing architecture, preferences
- Short and clear (one line each)
- Answerable without project file access

Return JSON only:
{
  "task_type": "feature|bugfix|testing|refactoring|documentation|other",
  "questions": ["Question 1?", "Question 2?", ...]
}`

	userPrompt := fmt.Sprintf(`Task: %s

Additional Details: %s

Analyze this task and generate context questions.`, task, details)

	response, err := c.SendMessage(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}

	var result AnalysisResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Ensure at least some questions
	if len(result.Questions) == 0 {
		result.Questions = []string{
			"What constraints or requirements should I be aware of?",
			"Are there any existing patterns or conventions to follow?",
		}
	}

	return &result, nil
}
