package main

import (
	"log"
	"os"

	"github.com/alecthomas/kingpin"
)

func main() {
	app := kingpin.New("db-index-task", "")
	app.Command("load", "")
	app.Command("split", "")
	migrationCmd := app.Command("migration", "help")
	migrationCmd.Command("up", "")
	migrationCmd.Command("down", "")

	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	switch cmd {

	}
}
