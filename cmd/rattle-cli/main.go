package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.UseShortOptionHandling = true
	app.Name = "rattle-bootstrap"
	app.Commands = append(app.Commands, newStartCommand())
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
