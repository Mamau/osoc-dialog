package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
	// Consider renaming to StmtBuilder or StatementBuilder to make it obvious.
	Builder squirrel.StatementBuilderType
}
type ProxyMysql struct {
	*DB
}

func Open(opts ...Option) (*DB, error) {
	o := &options{
		host:                 "localhost",
		port:                 3306,
		user:                 "root",
		parseTime:            true,
		verificationRequired: true,
		verificationTimeout:  5 * time.Second,
		////////////////////////////////////////////////////////////////
		// !!! Performance-related options !!!
		// Do not override them blindly.
		// Based on https://www.alexedwards.net/blog/configuring-sqldb.
		maxOpenConnections: 25,
		maxIdleConnections: 25,
		connMaxLifetime:    5 * time.Minute,
		///
		////////////////////////////////////////////////////////////////
	}
	for _, opt := range opts {
		opt(o)
	}

	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=%v",
		o.user,
		o.password,
		o.host,
		o.port,
		o.dbName,
		o.parseTime,
	)

	sqlxDB, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("mysql: open connection: %w", err)
	}

	db := DB{
		DB:      sqlxDB,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
	db.SetMaxOpenConns(o.maxOpenConnections)
	db.SetMaxIdleConns(o.maxIdleConnections)
	db.SetConnMaxLifetime(o.connMaxLifetime)

	if o.verificationRequired {
		ctx, cancel := context.WithTimeout(context.Background(), o.verificationTimeout)
		defer cancel()
		if err = db.PingContext(ctx); err != nil {
			return &db, fmt.Errorf("mysql: ping connection: %w", err)
		}
	}

	return &db, nil
}
