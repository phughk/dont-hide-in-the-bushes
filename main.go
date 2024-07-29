package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "filepath",
				Aliases:  []string{"f"},
				Required: false,
				Value:    "save-data.db",
				Usage:    "Filepath used for storing save data",
			}, &cli.IntFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Required: false,
				Value:    0,
				Usage:    "Port to listen on",
			},
		},
		Action: start_program,
	}
	app.Run(os.Args)
}

func start_program(context *cli.Context) error {
	network := &Network{}
	host, port, err := network.Bind(0)
	if err != nil {
		fmt.Printf("Failed to bind port: %v\n", err)
	} else {
		fmt.Printf("Listening on %s:%d\n", host, port)
	}
	defer func() {
		errs := network.Close()
		for err := range errs {
			fmt.Printf("Error closing connection: %v\n", err)
		}
	}()
	return nil
}
