package database

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

const DatabasePingTimeout = 10

type Database struct {
	Pool *pgxpool.Pool
	Log  zerolog.Logger
}
