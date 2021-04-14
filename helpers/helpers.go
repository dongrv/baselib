package sqlite

import (
	"baselib/db/common"
	"baselib/logger"
	"database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	Gorm *gorm.DB
)

func InitSqlite() *gorm.DB {
	var err error
	Gorm, err = gorm.Open("sqlite3", "./wss.db")
	if err != nil {
		if err != nil {
			logger.Sugar.Errorf("Failed to open Sqlite database. ")
			return Gorm
		}
	}
	logger.Sugar.Info("Sqlite connection succeeded")
	return Gorm
}

func Close() error {
	err := Gorm.Close()
	logger.Sugar.Info("Sqlite connection closed.")
	return err
}

func Insert(table string, insert map[string]interface{}) *gorm.DB {
	sql, values := common.InsertSetFormat(table, insert)
	return Gorm.Exec(sql, values...)
}

//更新
func Update(table string, where map[string]interface{}, update map[string]interface{}) *gorm.DB {
	sql, values := common.UpdateSetFormat(table, where, update)
	return Gorm.Exec(sql, values...)
}

//删除
func Delete(table string, where map[string]interface{}) *gorm.DB {
	sql := "delete from " + table
	whereSql, values := common.WhereSetFormat(where)
	sql += whereSql
	logger.Sugar.Info(sql)
	return Gorm.Exec(sql, values...)
}

// 获取一条数据
func FetchOne(table string, fields string, where map[string]interface{}) (row *gorm.DB) {
	sql := "select " + fields + " from " + table
	whereSql, values := common.WhereSetFormat(where)
	sql += whereSql
	logger.Sugar.Info(sql)
	row = Gorm.Raw(sql, values...)
	return
}

// 获取多条数据
func FetchAll(table string, fields string, where map[string]interface{}) (rows *sql.Rows) {
	var err error
	sql := "select " + fields + " from " + table
	whereSql, values := common.WhereSetFormat(where)
	sql += whereSql
	logger.Sugar.Info(sql)
	rows, err = Gorm.Raw(sql, values...).Rows()
	if err != nil {
		logger.Sugar.Errorf("Database query error. Error: %s", err)
	}
	return
}

// 关闭多条数据句柄
func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		logger.Sugar.Errorf("Failed to close result set. Error: %s", err)
	}
}
