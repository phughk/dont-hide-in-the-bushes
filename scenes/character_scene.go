package scenes

import (
	"fmt"
	"github.com/chzyer/readline"
)

type CharacterScene struct {
}

func (scene *CharacterScene) Prompt(cli *readline.Instance) {
	fmt.Println("This is the character scene")
}

func (scene *CharacterScene) Submit(cli *readline.Instance, input string) Scene {
	fmt.Printf("You wrote '%s'\n", input)
	return scene
}
