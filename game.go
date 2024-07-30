package main

import (
	"context"
	tl "github.com/phughk/termloop"
)

type Game struct {
	tlGame *tl.Game
}

func NewGame(ctx context.Context) (*Game, context.CancelFunc) {
	g := tl.NewGame()
	g.Screen().SetFps(60)
	ctx, cancel := context.WithCancel(ctx)
	game := &Game{tlGame: g}
	go g.StartCtx(ctx)
	return game, cancel
}
