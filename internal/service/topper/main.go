package topper

import (
	"encoding/csv"
	"io"
	"log"

	"github.com/PaulRaUnite/uni-db-index-task/internal/data"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

func Run(r io.Reader, limit int) error {
	unm, err := gocsv.NewUnmarshaller(csv.NewReader(r), data.Record{})
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i <= limit; i++ {
		rawRecord, err := unm.Read()
		if err == io.EOF {
			log.Println("EOF")
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "failed to decode record")
		}
		log.Println(rawRecord.(data.Record))
	}
	return nil
}
