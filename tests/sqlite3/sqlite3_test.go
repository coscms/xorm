package tests

import (
	"database/sql"
	"os"
	"testing"

	. "../"
	"github.com/coscms/xorm/core"
	"github.com/coscms/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func newSqlite3Engine() (*xorm.Engine, error) {
	os.Remove("./test.db")
	return xorm.NewEngine("sqlite3", "./test.db")
}

func newSqlite3DriverDB() (*sql.DB, error) {
	os.Remove("./test.db")
	return sql.Open("sqlite3", "./test.db")
}

func TestSqlite3(t *testing.T) {
	engine, err := newSqlite3Engine()
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
	BaseTestAllSnakeMapper(engine, t)
	BaseTestAll2(engine, t)
	BaseTestAll3(engine, t)
}

func TestSqlite3WithCache(t *testing.T) {
	engine, err := newSqlite3Engine()
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
	BaseTestAllSnakeMapper(engine, t)
	BaseTestAll2(engine, t)
}

func TestSqlite3SameMapper(t *testing.T) {
	engine, err := newSqlite3Engine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetMapper(core.SameMapper{})
	engine.ShowSQL = ShowTestSql
	engine.ShowErr = ShowTestSql
	engine.ShowWarn = ShowTestSql
	engine.ShowDebug = ShowTestSql

	BaseTestAll(engine, t)
	BaseTestAllSameMapper(engine, t)
	BaseTestAll2(engine, t)
	BaseTestAll3(engine, t)
}

func TestSqlite3WithCacheSameMapper(t *testing.T) {
	engine, err := newSqlite3Engine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetMapper(core.SameMapper{})
	engine.SetDefaultCacher(NewCacher())
	engine.ShowSQL = ShowTestSql
	engine.ShowErr = ShowTestSql
	engine.ShowWarn = ShowTestSql
	engine.ShowDebug = ShowTestSql

	BaseTestAll(engine, t)
	BaseTestAllSameMapper(engine, t)
	BaseTestAll2(engine, t)
}

const (
	createTableSqlite3 = "CREATE TABLE IF NOT EXISTS `big_struct` (`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, `name` TEXT NULL, `title` TEXT NULL, `age` TEXT NULL, `alias` TEXT NULL, `nick_name` TEXT NULL);"
	dropTableSqlite3   = "DROP TABLE IF EXISTS `big_struct`;"
)

func BenchmarkSqlite3DriverInsert(t *testing.B) {
	DoBenchDriver(newSqlite3DriverDB, createTableSqlite3, dropTableSqlite3,
		DoBenchDriverInsert, t)
}

func BenchmarkSqlite3DriverFind(t *testing.B) {
	DoBenchDriver(newSqlite3DriverDB, createTableSqlite3, dropTableSqlite3,
		DoBenchDriverFind, t)
}

func BenchmarkSqlite3NoCacheInsert(t *testing.B) {
	t.StopTimer()
	engine, err := newSqlite3Engine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchInsert(engine, t)
}

func BenchmarkSqlite3NoCacheFind(t *testing.B) {
	t.StopTimer()
	engine, err := newSqlite3Engine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchFind(engine, t)
}

func BenchmarkSqlite3NoCacheFindPtr(t *testing.B) {
	t.StopTimer()
	engine, err := newSqlite3Engine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchFindPtr(engine, t)
}

func BenchmarkSqlite3CacheInsert(t *testing.B) {
	t.StopTimer()
	engine, err := newSqlite3Engine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())
	DoBenchInsert(engine, t)
}

func BenchmarkSqlite3CacheFind(t *testing.B) {
	t.StopTimer()
	engine, err := newSqlite3Engine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())
	DoBenchFind(engine, t)
}

func BenchmarkSqlite3CacheFindPtr(t *testing.B) {
	t.StopTimer()
	engine, err := newSqlite3Engine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())
	DoBenchFindPtr(engine, t)
}
