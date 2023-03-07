package postgres

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vibin18/bse_shares/handler"
	"github.com/vibin18/bse_shares/model"
	"log"
	"time"
)

type postgresRepository struct {
	pgSql *sql.DB
}

const (
	maxOpenDbConn = 10
	maxIdleDbConn = 5
	maxDbLifetime = 5 * time.Minute
)

func newPgDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func testPgDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func NewPgSQL(dsn string) (handler.HandlerRepo, error) {
	dbConn := &postgresRepository{}
	d, err := newPgDb(dsn)
	if err != nil {
		panic(err)

	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbConn.pgSql = d

	err = testPgDB(dbConn.pgSql)
	if err != nil {
		panic(err)
	}
	log.Println("Connected to database")

	// This is where the database Repository is connected with the updater Repository
	// Note: The return is an interface UpdaterRepository

	return dbConn, nil
}

func (d *postgresRepository) GetAllPurchaseReport() (model.ShareReports, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares model.ShareReports

	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select purchased_shares.name, sum(purchased_shares.count) as pc, sum(purchased_shares.count * purchased_shares.price) as pprice from purchased_shares group by purchased_shares.name order by name `

	rows, err := d.pgSql.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r model.ShareReport
		err := rows.Scan(
			&r.Name,
			&r.Count,
			&r.Total,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil
}

func (d *postgresRepository) GetAllSalesReport() ([]model.ShareReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []model.ShareReport

	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select sold_shares.name,sum(sold_shares.count) as count_sum, sum(sold_shares.count * sold_shares.price) as total from sold_shares group by sold_shares.name order by name`

	rows, err := d.pgSql.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r model.ShareReport
		err := rows.Scan(
			&r.Name,
			&r.Count,
			&r.Total,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}

func (d *postgresRepository) InsertNewShare(res model.Share) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into share_names (name, id, created_at, updated_at) values ($1, $2, $3, $4)`
	_, err := d.pgSql.ExecContext(ctx, stmt,
		res.Name,
		res.Id,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (d *postgresRepository) GetShareByID(id int) (model.Share, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var share model.Share
	query := `select Name, id, created_at, updated_at from share_names where id = $2`
	row := d.pgSql.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&share.Name,
		&share.Id,
		&share.CreatedAt,
		&share.UpdatedAt,
	)
	if err != nil {
		return share, err
	}
	return share, nil
}

func (d *postgresRepository) GetAllShares() ([]model.Share, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []model.Share
	query := `select name, id from share_names`
	rows, err := d.pgSql.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r model.Share
		err := rows.Scan(
			&r.Name,
			&r.Id,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}

func (d *postgresRepository) BuyShare(res model.SellBuyShare) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into purchased_shares (name, id, count, price, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`
	_, err := d.pgSql.ExecContext(ctx, stmt,
		res.Name,
		res.Id,
		res.Count,
		res.Price,
		res.CreatedAt,
		res.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
func (d *postgresRepository) SellShare(res model.SellBuyShare) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into sold_shares (name, id, count, price, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`
	_, err := d.pgSql.ExecContext(ctx, stmt,
		res.Name,
		res.Id,
		res.Count,
		res.Price,
		res.CreatedAt,
		res.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (d *postgresRepository) GetAllSharesWithData() ([]model.TotalShare, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []model.TotalShare
	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select purchased_shares.id, purchased_shares.name, purchased_shares.count as pc, coalesce(sold_shares.count, 0) as sc from purchased_shares left outer join sold_shares on purchased_shares.id = sold_shares.id`

	rows, err := d.pgSql.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r model.TotalShare
		err := rows.Scan(
			&r.Id,
			&r.Name,
			&r.PCount,
			&r.SCount,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}

func (d *postgresRepository) GetAllPurchases() ([]model.SellBuyShare, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []model.SellBuyShare

	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select purchased_shares.id, purchased_shares.name, purchased_shares.count, purchased_shares.price, purchased_shares.updated_at, purchased_shares.created_at  from purchased_shares order by created_at`

	rows, err := d.pgSql.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r model.SellBuyShare
		err := rows.Scan(
			&r.Id,
			&r.Name,
			&r.Count,
			&r.Price,
			&r.UpdatedAt,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}

func (d *postgresRepository) GetAllSales() ([]model.SellBuyShare, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []model.SellBuyShare

	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select sold_shares.id, sold_shares.name, sold_shares.count, sold_shares.price, sold_shares.updated_at, sold_shares.created_at  from sold_shares order by created_at`

	rows, err := d.pgSql.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r model.SellBuyShare
		err := rows.Scan(
			&r.Id,
			&r.Name,
			&r.Count,
			&r.Price,
			&r.UpdatedAt,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}
