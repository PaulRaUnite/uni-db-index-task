package postgres

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func applySelector(query squirrel.SelectBuilder, selector data.GoodSelector) squirrel.SelectBuilder {
	if selector.Description != nil {
		query = query.Where("to_tsvector(description) @@ to_tsquery(?)", strings.Join(strings.Split(*selector.Description, " "), " & "))
	}
	query = selector.ApplyTo(query, "id")
	return query
}
func (q GoodQ) All(selector data.GoodSelector) ([]data.Good, error) {
	query := squirrel.Select("*").From("goods")

	query = applySelector(query, selector)
	var goods []data.Good
	err := q.db.Select(&goods, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select goods")
	}
	return goods, nil
}

func (q GoodQ) AllCount(description *string) (int64, error) {
	query := squirrel.Select("COUNT(*)").From("goods")
	if description != nil {
		query = query.Where("to_tsvector(description) @@ to_tsquery(?)", strings.Join(strings.Split(*description, " "), " & "))
	}

	var count int64
	err := q.db.GetContext(context.Background(), &count, query)
	if err != nil {
		return -1, errors.Wrap(err, "failed to select goods count")
	}
	return count, nil
}
