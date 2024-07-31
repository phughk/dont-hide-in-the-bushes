package scenes

import (
	"context"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"strings"
)

type CharacterScene struct {
	AlreadyPrompted bool
}

func (scene *CharacterScene) Prompt(ctx context.Context, cli *readline.Instance, llm *ollama.LLM) {
	if !scene.AlreadyPrompted {
		prompt_lines := []string{
			"You are a wise and mysterious dungeon master.",
			"You are talking to a player directly via chat.",
			"Respond to the player’s actions and inquiries in a way that maintains the atmosphere and narrative of the game.",
			"Do not break character and stay in immersion.",
			//"Keep your answers short, up to 100 words.",
			"Here is the scene to introduce the player to:",
			"You are on a date, your date has mysteriously run out of the building and into the bushes in front of the house. You are debating calling an Uber.",
		}
		prompt := strings.Join(prompt_lines, "\n")
		completion, err := llm.Call(ctx, prompt, llms.WithTemperature(0.8), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if chunk != nil {
				fmt.Print(string(chunk))
			}
			return nil
		}), llms.WithMaxLength(100))
		if err != nil {
			panic(err)
		}
		_ = completion
		scene.AlreadyPrompted = true
	}
	fmt.Println("Finished prompt")
}

func (scene *CharacterScene) Submit(ctx context.Context, cli *readline.Instance, llm *ollama.LLM, input string) Scene {
	//fmt.Printf("You wrote '%s'\n", input)
	completion, err := llm.Call(ctx, input, llms.WithTemperature(0.8), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if err != nil {
		panic(err)
	}
	_ = completion
	return scene
}
