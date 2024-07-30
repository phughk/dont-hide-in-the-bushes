package main

import (
	"context"
	"dont-hide-in-the-bushes/scenes"
	"github.com/chzyer/readline"
	"log"
)

type CliGame struct {
	cli *readline.Instance
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
	go func() {
		<-ctx.Done()
		err := cli.Close()
		if err != nil {
			log.Panicf("Failed to close terminal: %e", err)
		}
	}()
	cli.CaptureExitSignal()
	return &CliGame{
		cli: cli,
	}, err
}

func (g *CliGame) StartGame() {
	var nextScene scenes.Scene = &scenes.CharacterScene{}
	for nextScene != nil {
		nextScene = g.EnterScene(nextScene)
	}
}

func (g *CliGame) EnterScene(s scenes.Scene) scenes.Scene {
	s.Prompt(g.cli)
	line, err := g.cli.Readline()
	if err != nil {
		return nil
	}
	return s.Submit(g.cli, line)
}
