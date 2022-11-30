package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mauryasaurav/timescale_database/utils/helpers"
)

/*PostgresConnect - Connect function which is used in the main package for database connection */
func ConnectTimescaleDB() (*sqlx.DB, error) {

	/* Loading TOML file */
	config, err := helpers.LoadEnvFile()
	if err != nil {
		return nil, err
	}

	host := config.Get("timescale.host").(string)
	port := config.Get("timescale.port").(int64)
	user := config.Get("timescale.user").(string)
	password := config.Get("timescale.password").(string)
	dbName := config.Get("timescale.dbName").(string)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
