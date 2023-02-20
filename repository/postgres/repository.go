package postgres

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vibin18/bse_shares/updater"
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

func NewPgSQL(dsn string) (updater.ShareListUpdaterRepository, error) {
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

func (d *postgresRepository) GetAllPurchaseReport() (updater.ShareReports, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares updater.ShareReports

	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select purchased_shares.name, sum(purchased_shares.count) as pc, sum(purchased_shares.count * purchased_shares.price) as pprice from purchased_shares group by purchased_shares.name order by name `

	rows, err := d.pgSql.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r updater.ShareReport
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
