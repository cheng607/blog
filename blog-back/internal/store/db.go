package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	// 连接池优化
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	// 连接最大存活时间，避免服务端主动关闭后客户端仍持有无效连接
	db.SetConnMaxLifetime(5 * 60 * time.Second)
	return db, nil
}
