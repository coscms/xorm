package tests

import (
	"database/sql"
	"testing"

	. ".."
	"github.com/coscms/xorm"
	_ "github.com/ziutek/mymysql/godrv"
)

/*
CREATE DATABASE IF NOT EXISTS xorm_test CHARACTER SET
utf8 COLLATE utf8_general_ci;
*/

func TestMyMysql(t *testing.T) {
	err := mymysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}
	engine, err := xorm.NewEngine("mymysql", "xorm_test/root/")
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
	BaseTestAll3(engine, t)
}

func TestMyMysqlWithCache(t *testing.T) {
	err := mymysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}
	engine, err := xorm.NewEngine("mymysql", "xorm_test2/root/")
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

func newMyMysqlEngine() (*xorm.Engine, error) {
	return xorm.NewEngine("mymysql", "xorm_test2/root/")
}

func newMyMysqlDriverDB() (*sql.DB, error) {
	return sql.Open("mymysql", "xorm_test2/root/")
}

func BenchmarkMyMysqlDriverInsert(t *testing.B) {
	DoBenchDriver(newMyMysqlDriverDB, CreateTableMySql, DropTableMySql,
		DoBenchDriverInsert, t)
}

func BenchmarkMyMysqlDriverFind(t *testing.B) {
	DoBenchDriver(newMyMysqlDriverDB, CreateTableMySql, DropTableMySql,
		DoBenchDriverFind, t)
}

func mymysqlDdlImport() error {
	engine, err := xorm.NewEngine("mymysql", "/root/")
	if err != nil {
		return err
	}
	engine.ShowSQL = ShowTestSql
	engine.ShowErr = ShowTestSql
	engine.ShowWarn = ShowTestSql
	engine.ShowDebug = ShowTestSql

	sqlResults, _ := engine.Import("../testdata/mysql_ddl.sql")
	engine.LogDebug("sql results: %v", sqlResults)
	engine.Close()
	return nil
}

func BenchmarkMyMysqlNoCacheInsert(t *testing.B) {
	engine, err := newMyMysqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	DoBenchInsert(engine, t)
}

func BenchmarkMyMysqlNoCacheFind(t *testing.B) {
	engine, err := newMyMysqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	//engine.ShowSQL = true
	DoBenchFind(engine, t)
}

func BenchmarkMyMysqlNoCacheFindPtr(t *testing.B) {
	engine, err := newMyMysqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	//engine.ShowSQL = true
	DoBenchFindPtr(engine, t)
}

func BenchmarkMyMysqlCacheInsert(t *testing.B) {
	engine, err := newMyMysqlEngine()
	if err != nil {
		t.Error(err)
		return
	}

	defer engine.Close()
	engine.SetDefaultCacher(NewCacher())

	DoBenchInsert(engine, t)
}

func BenchmarkMyMysqlCacheFind(t *testing.B) {
	engine, err := newMyMysqlEngine()
	if err != nil {
		t.Error(err)
		return
	}

	defer engine.Close()
	engine.SetDefaultCacher(NewCacher())

	DoBenchFind(engine, t)
}

func BenchmarkMyMysqlCacheFindPtr(t *testing.B) {
	engine, err := newMyMysqlEngine()
	if err != nil {
		t.Error(err)
		return
	}

	defer engine.Close()
	engine.SetDefaultCacher(NewCacher())

	DoBenchFindPtr(engine, t)
}
