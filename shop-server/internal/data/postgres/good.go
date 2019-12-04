package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type GoodSelector struct {
}

// All return
func (q GoodQ) All(selector GoodSelector) ([]data.Good, error) {
	query := squirrel.Select("*").From("goods")

	var goods []data.Good
	err := q.db.Select(&goods, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select goods")
	}
	return goods, nil
}
