package orm

import (
	"baselib/db/common"
	"baselib/logger"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"time"
)

// 关系型数据库管理系统
// 数据持久化
// gorm 封装版

var (
	Gorm            *gorm.DB // 全局DB实例，当前包不可有重名变量
	dsn             string   // 连接数据源
	maxOpenConns    int      // 连接池最大数量，设置为0表示无限制
	maxIdleConns    int      // 空闲时最大连接数
	connMaxLifetime int64    // 连接最大过期时间
	maxIdleTime     int64    // 空闲时连接最大过期时间
	logMode         int      // 日志模式
)

type GormLogger struct{}

// 使用zap打印日志
func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		break // TODO 临时关闭
		logger.Sugar.Info(
			"sql",
			zap.String("type", "sql"),
			zap.Any("src", v[1]),
			//zap.Any("duration", v[2]),
			zap.Any("sql", v[3]),
			zap.Any("values", v[4]),
		//zap.Any("rows_returned", v[5]),
		)
	case "log":
		logger.Sugar.Error(
			"log",
			zap.Any("src", v[1]),
			zap.Any("gorm", v[2]),
		)
	}
}

// 连接设置
func setting() {
	dbConfig := struct {
		Host     string
		Port     string
		DBName   string
		Username string
		Password string
		Options  string
	}{
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.database"),
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.options"),
	}
	dsn = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Options,
	)
	logMode = viper.GetInt("mysql.logMode")
	maxOpenConns = viper.GetInt("mysql.maxOpenConns")
	maxIdleConns = viper.GetInt("mysql.maxIdleConns")
	connMaxLifetime = viper.GetInt64("mysql.connMaxLifetime")
	maxIdleTime = viper.GetInt64("mysql.maxIdleTime")
}

// 初始化连接
func InitMySQL() *gorm.DB {
	setting()
	var err error
	Gorm, err = gorm.Open("mysql", dsn)
	if err != nil {
		logger.Sugar.Errorf("Failed to connect to MySQL database. The dsn is `%s`. Error: %s", dsn, err)
		return Gorm
	}
	Gorm.DB().SetMaxOpenConns(maxOpenConns)
	Gorm.DB().SetMaxIdleConns(maxIdleConns)
	Gorm.DB().SetConnMaxLifetime(time.Duration(connMaxLifetime))
	Gorm.DB().SetConnMaxIdleTime(time.Duration(maxIdleTime))
	Gorm.SingularTable(false) // 表名启用复数形式
	if logMode == 1 {
		Gorm.LogMode(true) // 默认只打印错误日志;true 打印详细日志;false 不打印日志
	}
	Gorm.SetLogger(&GormLogger{})
	err = Gorm.DB().Ping()
	if err != nil {
		logger.Sugar.Errorf("Failed to ping MySQL server. Error: %s", err)
	}
	logger.Sugar.Info("MySQL connection succeeded")
	return Gorm
}

// 执行SQL源语
// 插入
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
func FetchOne(table string, fields []string, where map[string]interface{}) (row *gorm.DB) {
	fieldsStr := strings.Join(fields, ",")
	sql := "select " + fieldsStr + " from " + table
	whereSql, values := common.WhereSetFormat(where)
	sql += whereSql
	logger.Sugar.Info(sql)
	row = Gorm.Raw(sql, values...)
	return
}

// 获取多条数据
func FetchAll(table string, fields []string, where map[string]interface{}) (rows *sql.Rows) {
	var err error
	fieldsStr := strings.Join(fields, ",")
	sql := "select " + fieldsStr + " from " + table
	whereSql, values := common.WhereSetFormat(where)
	sql += whereSql
	logger.Sugar.Info(sql)
	rows, err = Gorm.Raw(sql, values...).Rows()
	if err != nil {
		logger.Sugar.Errorf("Database query error. Error: %s", err)
	}
	return
}

// 返回新插入id，插入的表需包含 自增主键
func RawInsertId(sqlStr string) (id int64, err error) {
	logger.Sugar.Info("raw insert return id sql:", sqlStr)
	var stmt *sql.Stmt
	var res sql.Result
	stmt, err = Gorm.DB().Prepare(sqlStr)
	res, err = stmt.Exec()
	id, err = res.LastInsertId()
	if err != nil {
		id = -1 // 插入报错
	}
	return
}

// 原生sql语句插入
func RawInsert(sqlStr string) (err error) {
	logger.Sugar.Info("raw insert sql:", sqlStr)
	var stmt *sql.Stmt
	stmt, err = Gorm.DB().Prepare(sqlStr)
	_, err = stmt.Exec()
	return
}

//原生sql更新
func RawUpdate(sql string) *gorm.DB {
	logger.Sugar.Info("raw update sql:", sql)
	return Gorm.Exec(sql)
}

//原生sql删除
func RawDelete(sql string) *gorm.DB {
	logger.Sugar.Info("raw delete sql:", sql)
	return Gorm.Exec(sql)
}

//原生sql获取多条
func RawFetchAll(sql string) (rows *sql.Rows) {
	logger.Sugar.Info("raw fetch all sql:", sql)
	start := time.Now()
	var err error
	rows, err = Gorm.Raw(sql).Rows()
	if err != nil {
		logger.Sugar.Errorf("Database query error. Error: %s", err)
	}
	if viper.GetInt("wss.debug") == 1 {
		secs := time.Since(start).Seconds() // 耗时
		if secs > viper.GetFloat64("wss.slowTimeLimit") {
			logger.Sugar.Infof("slow sql:%s,time consum:%fs", sql, secs)
		}
	}
	return
}

//原生sql获取一条数据
func RawFetchOne(sql string) (row *gorm.DB) {
	logger.Sugar.Info("raw fetch one sql:", sql)
	start := time.Now()
	row = Gorm.Raw(sql)
	if viper.GetInt("wss.debug") == 1 {
		secs := time.Since(start).Seconds() // 耗时
		if secs > viper.GetFloat64("wss.slowTimeLimit") {
			logger.Sugar.Infof("slow sql:%s,time consum:%fs", sql, secs)
		}
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

// Ping 检测
func Ping() error {
	return Gorm.DB().Ping()
}

// 关闭连接
func Close() error {
	err := Gorm.Close()
	logger.Sugar.Info("MySQL connection closed.")
	return err
}

func ScanAllRows(rows *sql.Rows, result interface{}) error {
	defer CloseRows(rows)

	// type1是*[]*Result
	type1 := reflect.TypeOf(result)
	if type1.Kind() != reflect.Ptr {
		return errors.New("第一个参数必须是指针")
	}
	// type2是[]*Result
	type2 := type1.Elem() // 解指针后的类型
	if type2.Kind() != reflect.Slice {
		return errors.New("第一个参数必须指向切片")
	}

	// type3是*Result
	type3 := type2.Elem()
	if type3.Kind() != reflect.Ptr {
		return errors.New("切片元素必须是指针类型")
	}

	// type4是Result
	type4 := type3.Elem()
	if type4.Kind() != reflect.Struct {
		return errors.New("切片元素必须是指针类型")
	}

	for rows.Next() {
		//  type3.Elem()是Result, elem是*Result
		elem := reflect.New(type3.Elem())
		// 传入*Result
		Gorm.ScanRows(rows, elem.Interface())
		// reflect.ValueOf(result).Elem()是[]*Result，Elem是*User，newSlice是[]*Result
		newSlice := reflect.Append(reflect.ValueOf(result).Elem(), elem)
		// 扩容后的slice赋值给*result
		// reflect.ValueOf(result).Elem()是[]Result
		reflect.ValueOf(result).Elem().Set(newSlice)
	}

	return nil
}
