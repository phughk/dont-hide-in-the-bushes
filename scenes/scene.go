package scenes

import "github.com/chzyer/readline"

type Scene interface {
	// Prompt Display the initial prompt upon entering the scene
	Prompt(cli *readline.Instance)
	// Submit The player is submitting an answer, if the scene returned is nil the game falls back to the previous scene or terminates
	// You can return the self-reference to stay in the scene (for example for multiple choice)
	Submit(cli *readline.Instance, input string) Scene
}
