package db

import (
	"database/sql"
	"fmt"

	"github.com/DenisGubenko/ideasymbols/models"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	// pq driver will be used in sql package
	_ "github.com/lib/pq"
)

type postgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(userName, password, dbName, host string,
	port int, sslMode string) (Storage, error) {
	var err error
	connection := new(postgresStorage)
	connection.db, err = sql.Open(
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%v sslmode=%s",
			userName, password, dbName, host, port, sslMode))

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return connection, nil
}

func (p *postgresStorage) CreateOrder(request *models.Order) error {
	tx, err := p.db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = tx.Exec(
		` INSERT INTO order_storage ( `+
			` content `+
			` ) VALUES ($1) ON CONFLICT (content) DO UPDATE SET active=$2, counter=$3`,
		request.Content, true, 0)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Wrap(err, errRollback.Error())
		}
		return errors.WithStack(err)
	}

	if err = tx.Commit(); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Wrap(err, errRollback.Error())
		}
		return errors.WithStack(err)
	}

	return nil
}

func (p *postgresStorage) GetRandomOrderContent() (*models.Order, error) {
	query := `WITH x AS (` +
		` SELECT id FROM order_storage WHERE active=$1 ORDER BY random() LIMIT 1 FOR UPDATE ` +
		` ) SELECT array_agg(id) FROM x;`
	result := p.db.QueryRow(query, true)

	var itemIDs []sql.NullInt64
	if err := result.Scan(pq.Array(&itemIDs)); err != nil {
		return nil, errors.WithStack(err)
	}

	query = ` SELECT content, active, counter FROM order_storage WHERE id=$1 `
	result = p.db.QueryRow(query, itemIDs[0].Int64)

	var order models.Order
	if err := result.Scan(&order.Content, &order.Active, &order.Counter); err != nil {
		return nil, errors.WithStack(err)
	}

	query = `UPDATE order_storage SET counter=$1 WHERE id=$2 `
	tx, err := p.db.Begin()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer stmt.Close()

	counter := order.Counter + 1
	if _, err = stmt.Exec(&counter, &itemIDs[0].Int64); err != nil {
		return nil, errors.WithStack(err)
	}

	if err = tx.Commit(); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return nil, errors.Wrap(err, errRollback.Error())
		}
		return nil, errors.WithStack(err)
	}

	return &order, nil
}

func (p *postgresStorage) InactiveRandomOrder() error {
	query := `WITH x AS (` +
		` SELECT id FROM order_storage WHERE active=true ORDER BY random() LIMIT 1 FOR UPDATE ` +
		` ) SELECT array_agg(id) FROM x;`
	result := p.db.QueryRow(query)

	var itemIDs []sql.NullInt64
	if err := result.Scan(pq.Array(&itemIDs)); err != nil {
		return errors.WithStack(err)
	}

	tx, err := p.db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	query = ` UPDATE order_storage SET active=$1 WHERE id=$2 `
	stmt, err := tx.Prepare(query)
	if err != nil {
		return errors.WithStack(err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(false, &itemIDs[0].Int64); err != nil {
		return errors.WithStack(err)
	}

	if err = tx.Commit(); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Wrap(err, errRollback.Error())
		}
		return errors.WithStack(err)
	}

	return nil
}

func (p *postgresStorage) GetStatisticsOrder() (*uint64, []*models.Order, error) {
	query := ` SELECT count(id) FROM order_storage `
	result := p.db.QueryRow(query)

	var count uint64
	if err := result.Scan(&count); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	query = ` SELECT content, active, counter FROM order_storage WHERE counter > $1 ORDER BY content ASC `
	rows, err := p.db.Query(query, 0)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	orders := make([]*models.Order, 0)
	for rows.Next() {
		order := models.Order{}
		err = rows.Scan(&order.Content, &order.Active, &order.Counter)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		orders = append(orders, &order)
	}

	return &count, orders, nil
}
