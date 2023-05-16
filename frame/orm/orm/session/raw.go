package session

import (
	"database/sql"
	"orm/log"
	"strings"
)

type Session struct {
	db *sql.DB // sql.Open() 方法连接数据库成功之后返回的指针
	// 拼接 SQL 语句和 SQL 语句中占位符的对应值
	sql     strings.Builder
	sqlVars []interface{}
}

func New(db *sql.DB) *Session {
	return &Session{
		db: db,
	}
}
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}
func (s *Session) DB() *sql.DB {
	return s.db
}

// Raw 改变这sql,sqlVars的变量值
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString("")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec raw sql with sqlVars
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRaw get a record from db
func (s *Session) QueryRaw() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRaws get a list records from db
func (s *Session) QueryRaws() (rows *sql.Rows, err error) {
	// 清空 (s *Session).sql 和 (s *Session).sqlVars 两个变量,可以复用
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return nil, err
}
