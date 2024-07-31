package main

import (
	"context"
	"dont-hide-in-the-bushes/scenes"
	"github.com/chzyer/readline"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
)

type CliGame struct {
	cli *readline.Instance
	llm *ollama.LLM
}

func newCliGame(ctx context.Context) (*CliGame, error) {
	cli, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	readline.CaptureExitSignal(func() {
		cancel()
		err := cli.Close()
		if err != nil {
			panic(err)
		}
	})
	cli.CaptureExitSignal()
	go func() {
		<-ctx.Done()
		err := cli.Close()
		if err != nil {
			log.Panicf("Failed to close terminal: %e", err)
		}
	}()
	cli.CaptureExitSignal()

	// Connect to the llm
	llm, err := ollama.New(ollama.WithModel("tinyllama"))
	if err != nil {
		panic(err)
	}
	return &CliGame{
		cli: cli,
		llm: llm,
	}, err
}

func (g *CliGame) StartGame(ctx context.Context) {
	var nextScene scenes.Scene = &scenes.CharacterScene{AlreadyPrompted: true}
	for nextScene != nil {
		nextScene = g.EnterScene(ctx, nextScene)
	}
}

func (g *CliGame) EnterScene(ctx context.Context, s scenes.Scene) scenes.Scene {
	s.Prompt(ctx, g.cli, g.llm)
	line, err := g.cli.Readline()
	if err != nil {
		return nil
	}
	return s.Submit(ctx, g.cli, g.llm, line)
}
