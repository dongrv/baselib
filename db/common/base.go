package common

import (
	"baselib/logger"
	"strings"
)

//格式化插入sql
func InsertSetFormat(table string, insert map[string]interface{}) (sql string, values []interface{}) {
	sql = "insert into " + table
	fields := "("
	val := " values("
	for key, value := range insert {
		fields += key + ","
		val += "?,"
		values = append(values, value)
	}
	fields = strings.Trim(fields, ",") + ")"
	val = strings.Trim(val, ",") + ")"
	sql += fields + val
	logger.Sugar.Info("INSERT SQL:" + sql)
	return
}

//格式化更新sql
func UpdateSetFormat(table string, where map[string]interface{}, update map[string]interface{}) (sql string, values []interface{}) {
	sql = "update " + table + " set "
	for key, value := range update {
		sql += key + "=?" + ","
		values = append(values, value)
	}
	sql = strings.Trim(sql, ",")
	sql += " where "
	and := ""
	for key, value := range where {
		sql += and + key + "=?"
		and = " and "
		values = append(values, value)
	}
	logger.Sugar.Info("UPDATE SQL:" + sql)
	return
}

//格式查询sql
func WhereSetFormat(where map[string]interface{}) (sql string, values []interface{}) {
	sql = " where "
	and := ""
	for key, value := range where {
		sql += and + key + "=?"
		and = " and "
		values = append(values, value)
	}
	return
}
