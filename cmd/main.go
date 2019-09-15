package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/PaulRaUnite/uni-db-index-task/internal/data"

	"github.com/gocarina/gocsv"

	"github.com/alecthomas/kingpin"
)

func main() {
	app := kingpin.New("db-index-task", "")
	topCmd := app.Command("top", "")
	limitArg := topCmd.Arg("limit", "").Required().Int()
	fileArg := topCmd.Arg("filename", "").Required().File()
	app.Command("load", "")
	app.Command("split", "")

	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	switch cmd {
	case topCmd.FullCommand():
		unm, err := gocsv.NewUnmarshaller(csv.NewReader(*fileArg), data.Record{})
		if err != nil {
			log.Fatalln(err)
		}
		for i := 0; i <= *limitArg; i++ {
			rawRecord, err := unm.Read()
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(rawRecord.(data.Record))
		}
	}
}
