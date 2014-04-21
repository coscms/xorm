package tests

//
// +build windows

import (
	"database/sql"
	"testing"

	. ".."
	"github.com/coscms/xorm"
	_ "github.com/coscms/xorm/driver/godbc"
)

const mssqlConnStr = "driver={SQL Server};Server=192.168.20.135;Database=xorm_test; uid=sa; pwd=1234;"

func newMssqlEngine() (*xorm.Engine, error) {
	return xorm.NewEngine("odbc", mssqlConnStr)
}

func TestMssql(t *testing.T) {
	engine, err := newMssqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.ShowSQL = ShowTestSql
	engine.ShowErr = ShowTestSql
	engine.ShowWarn = ShowTestSql
	engine.ShowDebug = ShowTestSql

	BaseTestAll(engine, t)
	BaseTestAll2(engine, t)
}

func TestMssqlWithCache(t *testing.T) {
	engine, err := newMssqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())
	engine.ShowSQL = ShowTestSql
	engine.ShowErr = ShowTestSql
	engine.ShowWarn = ShowTestSql
	engine.ShowDebug = ShowTestSql

	BaseTestAll(engine, t)
	BaseTestAll2(engine, t)
}

func newMssqlDriverDB() (*sql.DB, error) {
	return sql.Open("odbc", mssqlConnStr)
}

const (
	createTableMssql = `IF NOT EXISTS (SELECT [name] FROM sys.tables WHERE [name] = 'big_struct' ) CREATE TABLE
		"big_struct" ("id" BIGINT PRIMARY KEY IDENTITY NOT NULL, "name" VARCHAR(255) NULL, "title" VARCHAR(255) NULL,
		"age" VARCHAR(255) NULL, "alias" VARCHAR(255) NULL, "nick_name" VARCHAR(255) NULL);
		`

	dropTableMssql = "IF EXISTS (SELECT * FROM sysobjects WHERE id = object_id(N'big_struct') and OBJECTPROPERTY(id, N'IsUserTable') = 1) DROP TABLE IF EXISTS `big_struct`;"
)

func BenchmarkMssqlDriverInsert(t *testing.B) {
	DoBenchDriver(newMssqlDriverDB, createTableMssql, dropTableMssql,
		DoBenchDriverInsert, t)
}

func BenchmarkMssqlDriverFind(t *testing.B) {
	DoBenchDriver(newMssqlDriverDB, createTableMssql, dropTableMssql,
		DoBenchDriverFind, t)
}

func BenchmarkMssqlNoCacheInsert(t *testing.B) {
	engine, err := newMssqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchInsert(engine, t)
}

func BenchmarkMssqlNoCacheFind(t *testing.B) {
	engine, err := newMssqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchFind(engine, t)
}

func BenchmarkMssqlNoCacheFindPtr(t *testing.B) {
	engine, err := newMssqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchFindPtr(engine, t)
}

func BenchmarkMssqlCacheInsert(t *testing.B) {
	engine, err := newMssqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())

	DoBenchInsert(engine, t)
}

func BenchmarkMssqlCacheFind(t *testing.B) {
	engine, err := newMssqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())

	DoBenchFind(engine, t)
}

func BenchmarkMssqlCacheFindPtr(t *testing.B) {
	engine, err := newMssqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())

	DoBenchFindPtr(engine, t)
}
