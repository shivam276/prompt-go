package enhancer

import (
	"strings"

	"promptgo/internal/templates"
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

const tip = "This prompt forces AI to think before coding. You stay in control of architecture decisions. The secret word prevents premature implementation. Your sketch keeps your mental model in the code."

// Enhance takes user input and returns an enhanced prompt with the system template
func Enhance(input Input) Output {
	// Replace the placeholder with the actual secret word
	systemPrompt := strings.ReplaceAll(templates.SystemPrompt, "{{SECRET_WORD}}", input.SecretWord)

	// Build the final prompt
	var builder strings.Builder
	builder.WriteString(systemPrompt)
	builder.WriteString("\n\n---\n\n")
	builder.WriteString("## MY TASK\n\n")
	builder.WriteString(input.Task)

	if input.Details != "" {
		builder.WriteString("\n\n## ADDITIONAL CONTEXT\n\n")
		builder.WriteString(input.Details)
	}

	return Output{
		EnhancedPrompt: builder.String(),
		Tip:            tip,
	}
}
