package validator

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"

	"github.com/pkg/errors"

	"github.com/PaulRaUnite/uni-db-index-task/internal/data"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gocarina/gocsv"
)

func Run(source io.Reader) error {
	errs := validation.Errors{}

	unm, err := gocsv.NewUnmarshaller(csv.NewReader(source), data.Record{})
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
