package postgres

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// All return
func (q CustomerQ) All(selector data.CustomerSelector) ([]data.Customer, error) {
	query := squirrel.Select("*").From("customers")
	if selector.Search != nil {
		query = query.Where("name % ?", *selector.Search).OrderBy(fmt.Sprintf("name <-> '%s' DESC", *selector.Search))
	}

	var customers []data.Customer
	err := q.db.Select(&customers, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select customers")
	}
	return customers, nil
}
