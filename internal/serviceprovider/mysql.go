package serviceprovider

import (
	"osoc-dialog/internal/config"

	"osoc-dialog/pkg/mysql"
)

func NewMysql(conf config.Mysql) (*mysql.DB, func(), error) {
	db, err := mysql.Open(
		mysql.Host(conf.Host),
		mysql.Port(conf.Port),
		mysql.User(conf.User),
		mysql.Password(conf.Password),
		mysql.DBName(conf.DbName),
		mysql.ParseTime(conf.ParseTime),
		mysql.MaxIdleConnections(25),
		mysql.VerificationRequired(false),
	)
	if err != nil {
		return nil, nil, err
	}

	closeDB := func() {
		_ = db.Close()
	}

	return db, closeDB, nil
}

func NewProxyMysql(conf config.ProxyMysql) (*mysql.ProxyMysql, func(), error) {
	db, err := mysql.Open(
		mysql.Host(conf.Host),
		mysql.Port(conf.Port),
		mysql.User(conf.User),
		mysql.Password(conf.Password),
		mysql.DBName(conf.DbName),
		mysql.ParseTime(conf.ParseTime),
		mysql.MaxIdleConnections(25),
		mysql.VerificationRequired(false),
	)
	if err != nil {
		return nil, nil, err
	}

	closeDB := func() {
		_ = db.Close()
	}

	return &mysql.ProxyMysql{DB: db}, closeDB, nil
}
