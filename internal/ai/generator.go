package ai

import (
	"context"
	"fmt"
	"strings"
)

type PromptRequest struct {
	Task       string
	Details    string
	TaskType   TaskType
	QA         map[string]string // Question -> Answer
	SecretWord string
}

// GeneratePrompt generates a comprehensive, task-specific prompt
func (c *Client) GeneratePrompt(ctx context.Context, req PromptRequest) (string, error) {
	systemPrompt := `You are an expert prompt engineer for software development.

Generate a comprehensive, structured prompt that guides a developer through implementing their task.

The prompt should:
1. Be specific to the task type and context provided
2. Follow a phase-based methodology:
   - PHASE 1: UNDERSTAND - Deep analysis, propose approaches, ask clarifying questions
   - PHASE 2: ALIGN - Design review, sketch requirements, implementation planning
   - PHASE 3: BUILD - Implementation (gated by secret word: "{{SECRET_WORD}}")
3. Incorporate the context from the Q&A
4. Include task-specific best practices
5. Suggest testing strategies appropriate for the task
6. Be practical and actionable

Do not use generic templates. Create a fully custom prompt tailored to THIS specific task.`

	// Build context section from Q&A
	qaContext := ""
	if len(req.QA) > 0 {
		qaContext = "\n\nContext from your answers:\n"
		for q, a := range req.QA {
			qaContext += fmt.Sprintf("- %s → %s\n", q, a)
		}
	}

	userPrompt := fmt.Sprintf(`Task Type: %s
Task: %s
Details: %s%s
Secret Word: %s

Generate the enhanced prompt now.`, req.TaskType, req.Task, req.Details, qaContext, req.SecretWord)

	response, err := c.SendMessage(ctx, systemPrompt, userPrompt)
	if err != nil {
		return "", fmt.Errorf("prompt generation failed: %w", err)
	}

	// Replace the placeholder with the actual secret word in the response
	response = strings.ReplaceAll(response, "{{SECRET_WORD}}", req.SecretWord)

	return response, nil
}
