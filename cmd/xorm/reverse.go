package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/admpub/log"
	"github.com/coscms/xorm"
	"github.com/coscms/xorm/core"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/ziutek/mymysql/godrv"
)

var CmdReverse = &Command{
	UsageLine: "reverse [-s] driverName datasourceName tmplPath [generatedPath]",
	Short:     "reverse a db to codes",
	Long: `
according database's tables and columns to generate codes for Go, C++ and etc.

    -s                Generated one go file for every table
    driverName        Database driver name, now supported four: mysql mymysql sqlite3 postgres
    datasourceName    Database connection uri, for detail infomation please visit driver's project page
    tmplPath        Template dir for generated. the default templates dir has provide 1 template
    generatedPath    This parameter is optional, if blank, the default value is model, then will
                    generated all codes in model dir
`,
}

func init() {
	CmdReverse.Run = runReverse
	CmdReverse.Flags = map[string]bool{
		"-s": false,
		"-l": false,
	}
	log.Sync(true)
}

var (
	genJson bool = false
)

func printReversePrompt(flag string) {
}

type Tmpl struct {
	Tables  []*core.Table
	Imports map[string]string
	Model   string
}

func dirExists(dir string) bool {
	d, e := os.Stat(dir)
	switch {
	case e != nil:
		return false
	case !d.IsDir():
		return false
	}

	return true
}

func runReverse(cmd *Command, args []string) {
	num := checkFlags(cmd.Flags, args, printReversePrompt)
	if num == -1 {
		return
	}
	args = args[num:]

	if len(args) < 3 {
		fmt.Println("params error, please see xorm help reverse")
		return
	}

	isMultiFile := true
	if use, ok := cmd.Flags["-s"]; ok {
		isMultiFile = !use
	}

	curPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	var genDir string
	var model string
	if len(args) == 4 {
		genDir, err = filepath.Abs(args[3])
		if err != nil {
			fmt.Println(err)
			return
		}

		model = filepath.Base(genDir)
	} else {
		model = "model"
		genDir = filepath.Join(curPath, model)
	}

	dir, err := filepath.Abs(args[2])
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	if !dirExists(dir) {
		log.Errorf("Template %v path is not exist", dir)
		return
	}

	var langTmpl LangTmpl
	var ok bool
	lang := "go"
	var prefix string
	var ignore *regexp.Regexp

	cfgPath := filepath.Join(dir, "config")
	info, err := os.Stat(cfgPath)
	var configs map[string]string
	if err == nil && !info.IsDir() {
		configs = loadConfig(cfgPath)
		if l, ok := configs["lang"]; ok {
			lang = l
		}
		if j, ok := configs["genJson"]; ok {
			genJson, err = strconv.ParseBool(j)
		}

		//[SWH|+]
		if j, ok := configs["prefix"]; ok {
			prefix = j
		}
		if j, ok := configs["ignore"]; ok {
			ignore = regexp.MustCompile(j)
		}
	}
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	if langTmpl, ok = langTmpls[lang]; !ok {
		fmt.Println("Unsupported programing language", lang)
		return
	}

	os.MkdirAll(genDir, os.ModePerm)

	Orm, err := xorm.NewEngine(args[0], args[1])
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	tables, err := Orm.DBMetas()
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	structTmplCount := 0
	err = filepath.Walk(dir, func(f string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if info.Name() == "config" {
			return nil
		}

		bs, err := ioutil.ReadFile(f)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}

		t := template.New(f)
		t.Funcs(langTmpl.Funcs)

		tmpl, err := t.Parse(string(bs))
		if err != nil {
			log.Errorf("%v", err)
			return err
		}

		var w *os.File
		fileName := info.Name()
		newFileName := fileName[:len(fileName)-4]
		ext := filepath.Ext(newFileName)

		if !isMultiFile {
			w, err = os.Create(filepath.Join(genDir, newFileName))
			if err != nil {
				log.Errorf("%v", err)
				return err
			}

			imports := langTmpl.GenImports(tables)

			var tbls []*core.Table
			for _, table := range tables {
				//[SWH|+]
				if prefix != "" {
					table.Name = strings.TrimPrefix(table.Name, prefix)
				}
				if ignore != nil && ignore.MatchString(table.Name) {
					continue
				}
				tbls = append(tbls, table)
			}

			newbytes := bytes.NewBufferString("")

			t := &Tmpl{Tables: tbls, Imports: imports, Model: model}
			err = tmpl.Execute(newbytes, t)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}

			tplcontent, err := ioutil.ReadAll(newbytes)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}
			var source string
			if langTmpl.Formater != nil {
				source, err = langTmpl.Formater(string(tplcontent))
				if err != nil {
					log.Errorf("%v", err)
					return err
				}
			} else {
				source = string(tplcontent)
			}

			w.WriteString(source)
			w.Close()
			structTmplCount++
		} else {
			for _, table := range tables {
				//[SWH|+]
				if prefix != "" {
					table.Name = strings.TrimPrefix(table.Name, prefix)
				}
				if ignore != nil && ignore.MatchString(table.Name) {
					continue
				}
				// imports
				tbs := []*core.Table{table}
				imports := langTmpl.GenImports(tbs)

				w, err := os.Create(filepath.Join(genDir, unTitle(mapper.Table2Obj(table.Name))+ext))
				if err != nil {
					log.Errorf("%v", err)
					return err
				}

				newbytes := bytes.NewBufferString("")

				t := &Tmpl{Tables: tbs, Imports: imports, Model: model}
				err = tmpl.Execute(newbytes, t)
				if err != nil {
					log.Errorf("%v", err)
					return err
				}

				tplcontent, err := ioutil.ReadAll(newbytes)
				if err != nil {
					log.Errorf("%v", err)
					return err
				}
				var source string
				if langTmpl.Formater != nil {
					source, err = langTmpl.Formater(string(tplcontent))
					if err != nil {
						log.Errorf("%v-%v", err, string(tplcontent))
						return err
					}
				} else {
					source = string(tplcontent)
				}

				w.WriteString(source)
				w.Close()
				structTmplCount++
			}
		}

		return nil
	})
	if err != nil {
		log.Error(err)
	}
	if structTmplCount == 0 {
		log.Infof(`struct template not found.`)
	}
}
