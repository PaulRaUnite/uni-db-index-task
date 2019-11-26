package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"runtime/pprof"

	"github.com/PaulRaUnite/uni-db-index-task/db/internal/service/validator"

	"github.com/PaulRaUnite/uni-db-index-task/db/internal/config"
	"gitlab.com/distributed_lab/kit/kv"

	"github.com/PaulRaUnite/uni-db-index-task/db/internal/service/loader"

	"github.com/PaulRaUnite/uni-db-index-task/db/internal/service/splitter"
	"github.com/PaulRaUnite/uni-db-index-task/db/internal/service/topper"

	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
)

func main() {
	app := kingpin.New("db-index-task", "")
	cpuprofile := app.Flag("cpu-profile", "write cpu profile to `file`").String()
	topCmd := app.Command("top", "")

	limitArg := topCmd.Arg("limit", "").Required().Int()
	fileArg := topCmd.Arg("filename", "").Required().File()

	loadCmd := app.Command("load", "")
	loadFileArg := loadCmd.Arg("file", "").Required().File()

	splitCmd := app.Command("split", "")
	splitFrom := splitCmd.Arg("from", "").Required().File()
	splitInto := splitCmd.Arg("to", "").Required().String()

	validateCmd := app.Command("validate", "")
	validateFileArg := validateCmd.Arg("dataset", "").Required().File()

	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
	if cpuprofile != nil {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	switch cmd {
	case topCmd.FullCommand():
		defer (*fileArg).Close()
		err := topper.Run(*fileArg, *limitArg)
		if err != nil {
			log.Fatalln(errors.Wrap(err, "failed to print top records"))
		}
	case splitCmd.FullCommand():
		splitFrom := *splitFrom
		defer splitFrom.Close()
		lineCount, err := lineCounter(splitFrom)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = splitFrom.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatalln(err)
		}
		err = splitter.Run(splitFrom, splitFrom.Name(), *splitInto, lineCount)
		if err != nil {
			log.Fatalln(errors.Wrap(err, "failed to split dataset"))
		}
	case loadCmd.FullCommand():
		cfg := config.NewConfig(kv.MustFromEnv())
		defer (*loadFileArg).Close()
		err = loader.Run(cfg, *loadFileArg)

		if err != nil {
			log.Fatalln(errors.Wrap(err, "failed to load records to database"))
		}
	case validateCmd.FullCommand():
		file := *validateFileArg
		defer file.Close()
		err = validator.Run(file)
		if err != nil {
			log.Fatalln(errors.Wrap(err, "validation failed"))
		}
	}
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
