package main

import (
	"context"
	tl "github.com/phughk/termloop"
)

type FullTerminalGame struct {
	tlGame *tl.Game
}

func NewGame(ctx context.Context) (*FullTerminalGame, context.CancelFunc) {
	g := tl.NewGame()
	g.Screen().SetFps(60)
	ctx, cancel := context.WithCancel(ctx)
	game := &FullTerminalGame{tlGame: g}
	go g.StartCtx(ctx)
	return game, cancel
}
