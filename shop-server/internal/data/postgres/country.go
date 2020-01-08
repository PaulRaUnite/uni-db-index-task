package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (q CountryQ) All() ([]data.Country, error) {
	query := squirrel.Select("*").From("countries")

	var countries []data.Country
	err := q.db.Select(&countries, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select goods")
	}
	return countries, nil
}
