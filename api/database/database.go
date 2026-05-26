package database

import (
	"sync"
	"time"

	"github.com/TDiblik/project-template/api/utils"
	_ "github.com/jackc/pgx/v5/stdlib" // Register pgx with database/sql
	"github.com/jmoiron/sqlx"
)

var (
	db    *sqlx.DB
	once  sync.Once
	dbErr error
)

// CreateConnection initializes or returns the global DB pool.
func CreateConnection() (*sqlx.DB, error) {
	once.Do(func() {
		d, err := sqlx.Connect("pgx", utils.EnvData.DB_CONNECTION_STRING)
		if err != nil {
			dbErr = err
			return
		}

		// Pool tuning (adjust based on workload & DB limits)
		d.SetMaxOpenConns(50)                  // Max total connections
		d.SetMaxIdleConns(10)                  // Keep some idle for reuse
		d.SetConnMaxIdleTime(5 * time.Minute)  // Kill idle conns after 5 min
		d.SetConnMaxLifetime(30 * time.Minute) // Recycle conns after 30 min

		db = d
	})

	return db, dbErr
}
