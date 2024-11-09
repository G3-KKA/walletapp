package dbsink

import (
	"context"
	"fmt"
	"sync"
	"walletapp/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type pgxSink struct {
	pool *pgxpool.Pool
	cfg  config.Database

	l        *zerolog.Logger
	shutdown sync.Once
}

// New is a database wrapper ( pgxSink ) constructor.
//
// New doesnt connect to the database, only prepares everything.
//
// Constructor assumes that configuration is valid,
// internals of driver may or may not check config,
// so it is a caller responsibility to provide valid config.
func New(l *zerolog.Logger, cfg config.Database) *pgxSink {

	return &pgxSink{
		pool:     nil,
		cfg:      cfg,
		l:        l,
		shutdown: sync.Once{},
	}

}

// Ping goes through stages, any of it may error:
//   - parses config,
//   - assembles connection pool,
//   - pings database.
func (sink *pgxSink) Ping(ctx context.Context) error {

	dsn := buildDSN(sink.cfg.DBRole,
		sink.cfg.DBPassword,
		sink.cfg.DBAddress,
		sink.cfg.DBPort,
		sink.cfg.DBName,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {

		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)

	if err != nil {

		return err
	}

	err = pool.Ping(ctx)
	if err != nil {
		// Retry policy?
		pool.Close()

		return err
	}

	sink.pool = pool
	sink.l.Info().Msg("database sink is ready to use")

	return nil
}

// Shutdown closes all postgres connections,
// connection pool does not provides any errors,
// so if any problems even happens this function still return nil error.
//
// # Safe to multi-call.
func (sink *pgxSink) Shutdown(ctx context.Context) error {

	sink.shutdown.Do(func() {
		sink.pool.Close()
	})

	return nil
}

//nolint:all // Internals.
func buildDSN(role, rolepass, dbaddr, port, dbname string) string {

	// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		role,
		rolepass,
		dbaddr,
		port,
		dbname,
	)
}
