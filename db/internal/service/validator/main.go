package validator

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"

	data2 "github.com/PaulRaUnite/uni-db-index-task/db/internal/data"

	"github.com/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gocarina/gocsv"
)

func Run(source io.Reader) error {
	errs := validation.Errors{}

	unm, err := gocsv.NewUnmarshaller(csv.NewReader(source), data2.Record{})
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; ; i++ {
		_, m, err := unm.ReadUnmatched()
		if err == io.EOF {
			break
		}
		if err != nil {
			errs[strconv.Itoa(i)] = errors.Wrapf(err, "%v", m)
		}
	}
	return errs.Filter()
}
