package scenes

import (
	"context"
	"github.com/chzyer/readline"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Scene interface {
	// Prompt Display the initial prompt upon entering the scene
	Prompt(ctx context.Context, cli *readline.Instance, llm *ollama.LLM)
	// Submit The player is submitting an answer, if the scene returned is nil the game falls back to the previous scene or terminates
	// You can return the self-reference to stay in the scene (for example for multiple choice)
	Submit(ctx context.Context, cli *readline.Instance, llm *ollama.LLM, input string) Scene
}
