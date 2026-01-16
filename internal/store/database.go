package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	// implicit import for driver registration
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		// https://pkg.go.dev/fmt#Errorf
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Add enhanced configuration to the connection pool settings with:
	// db.SetMaxOpenConns(), db.SetMaxIdleConns(), and db.SetConnMaxIdleTime()
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}
	fmt.Println("Database opened successfully")
	return db, nil
}

/**
* 将 migrations 目录设置为 migrationsFs，并返回迁移错误。
* 迁移完成后，将 migrations 目录设置为 nil。
* 这个函数的主要作用是让数据库迁移能够支持 Go 的 embed 特性，
* 使得应用启动时可以自动从编译好的二进制文件中读取 SQL 脚本并更新数据库结构。
* 返回迁移错误。
* 通过这个设置，你的程序变成了一个独立的“全能包”。
* 无论你把这个编译好的程序扔到哪台机器上运行，
* 它都能自己从自己身体里掏出 SQL 脚本来升级数据库，完全不需要依赖外部的 .sql 文件。
* 这个功能在开发和部署时都非常方便，尤其是当你需要将应用打包成一个独立的二进制文件时。
* @param db *sql.DB 数据库连接
* @param migrationsFs fs.FS 一个实现了 Go标准库 fs.FS 接口的文件系统对象。这通常用于 Go 的 embed 功能，将 SQL 迁移文件直接打包进编译后的二进制文件中
* @param dir string 迁移目录
* @return error 迁移错误
 */

func MigrateFS(db *sql.DB, migrationsFs fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFs)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, migrationsDir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: set dialect: %w", err)
	}
	err = goose.Up(db, migrationsDir)
	if err != nil {
		return fmt.Errorf("migrate: up: %w", err)
	}
	return nil
}
