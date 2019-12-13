package names

import (
	"bufio"
	"io"

	"github.com/Masterminds/squirrel"
	"github.com/PaulRaUnite/uni-db-index-task/db/internal/config"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func Run(cfg config.Config, r io.Reader) error {
	s := bufio.NewScanner(r)
	db := cfg.DB()
	for s.Scan() {
		if err := s.Err(); err != nil {
			return err
		}
		name := s.Text()
		query := squirrel.Update("customers").Where("id = (select id from customers where name LIKE '' LIMIT 1)").
			SetMap(map[string]interface{}{"name": name})
		err := db.Exec(query)
		if err != nil {
			return errors.Wrap(err, "failed to insert invoice part")
		}
	}
	return nil
}
