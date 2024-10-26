package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(driver, host string, port uint16, dbName, user, password string, ssl bool) (*sqlx.DB, error) {
	var prefix string
	if ssl {
		prefix += "en"
	} else {
		prefix += "dis"
	}
	return sqlx.Open(driver,
		fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%sable",
			host, port, dbName, user, password, prefix,
		),
	)
	/*db, err := sqlx.Open(driver, dsn)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create database connection")

		return nil, err
	}

	// need to uncomment for homework-4
	// if err = db.Ping(); err != nil {
	//  log.Error().Err(err).Msgf("failed ping the database")

	//  return nil, err
	// }

	return db, nil*/
}
