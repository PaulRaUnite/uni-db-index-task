package postgres

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// All return
func (q UserQ) All(selector data.UserSelector) ([]data.User, error) {
	query := squirrel.Select("*").From("users")
	if selector.Name != nil {
		query = query.Where("name % ?", *selector.Name).OrderBy(fmt.Sprintf("name <-> '%s' DESC", *selector.Name))
	}

	var customers []data.User
	err := q.db.Select(&customers, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select customers")
	}
	return customers, nil
}
