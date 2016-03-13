package xorm

import (
	"fmt"

	"github.com/coscms/xorm/core"
)

var DefaultShowLog = map[string]bool{
	"sql":   true,
	"etime": true,
	"cache": true,
	"event": true,
	"base":  true,
	"other": true,
}

var defaultLogProcessor = func(tag string, format string, args []interface{}) (string, []interface{}) {
	if format == "" {
		if len(args) > 0 {
			args[0] = fmt.Sprintf("[%s] %v", tag, args[0])
		}
		return format, args
	}
	format = "[" + tag + "] " + format
	return format, args
}
var LogTagProcessor = map[string]func(tag string, format string, args []interface{}) (string, []interface{}){
	"cache": defaultLogProcessor,
	"event": defaultLogProcessor,
	"sql":   defaultLogProcessor,
	"etime": defaultLogProcessor,
	"base":  defaultLogProcessor,
	"other": defaultLogProcessor,
}

// =====================================
// 增加Engine结构体中的方法
// @author Admpub <swh@admpub.com>
// =====================================

func (engine *Engine) Init() {
	engine.showLog = DefaultShowLog
}

func (engine *Engine) SetTblMapper(mapper core.IMapper) {
	if prefixMapper, ok := mapper.(core.PrefixMapper); ok {
		engine.TablePrefix = prefixMapper.Prefix
	} else if suffixMapper, ok := mapper.(core.SuffixMapper); ok {
		engine.TableSuffix = suffixMapper.Suffix
	}
	engine.TableMapper = mapper
}

func (engine *Engine) OpenLog(types ...string) {
	if len(types) < 1 {
		for typ, _ := range engine.showLog {
			engine.showLog[typ] = true
			if typ == "sql" {
				engine.ShowSQL(engine.showLog[typ])
			}
		}
		return
	}
	for _, typ := range types {
		engine.showLog[typ] = true
		if typ == "sql" {
			engine.ShowSQL(engine.showLog[typ])
		}
	}
}

func (engine *Engine) CloseLog(types ...string) {
	if len(types) < 1 {
		for typ, _ := range engine.showLog {
			engine.showLog[typ] = false
			if typ == "sql" {
				engine.ShowSQL(engine.showLog[typ])
			}
		}
		return
	}
	for _, typ := range types {
		engine.showLog[typ] = false
		if typ == "sql" {
			engine.ShowSQL(engine.showLog[typ])
		}
	}
}

func (engine *Engine) TagLogError(tag string, contents ...interface{}) {
	if enable, _ := engine.showLog[tag]; !enable {
		return
	}
	if fn, ok := LogTagProcessor[tag]; ok {
		_, contents = fn(tag, "", contents)
		if contents == nil {
			return
		}
	}
	engine.LogError(contents...)
}

func (engine *Engine) TagLogErrorf(tag string, format string, contents ...interface{}) {
	if enable, _ := engine.showLog[tag]; !enable {
		return
	}
	if fn, ok := LogTagProcessor[tag]; ok {
		format, contents = fn(tag, format, contents)
		if contents == nil {
			return
		}
	}
	engine.LogErrorf(format, contents...)
}

// logging info
func (engine *Engine) TagLogInfo(tag string, contents ...interface{}) {
	if enable, _ := engine.showLog[tag]; !enable {
		return
	}
	if fn, ok := LogTagProcessor[tag]; ok {
		_, contents = fn(tag, "", contents)
		if contents == nil {
			return
		}
	}
	engine.LogInfo(contents...)
}

func (engine *Engine) TagLogInfof(tag string, format string, contents ...interface{}) {
	if enable, _ := engine.showLog[tag]; !enable {
		return
	}
	if fn, ok := LogTagProcessor[tag]; ok {
		format, contents = fn(tag, format, contents)
		if contents == nil {
			return
		}
	}
	engine.LogInfof(format, contents...)
}

// logging debug
func (engine *Engine) TagLogDebug(tag string, contents ...interface{}) {
	if enable, _ := engine.showLog[tag]; !enable {
		return
	}
	if fn, ok := LogTagProcessor[tag]; ok {
		_, contents = fn(tag, "", contents)
		if contents == nil {
			return
		}
	}
	engine.LogDebug(contents...)
}

func (engine *Engine) TagLogDebugf(tag string, format string, contents ...interface{}) {
	if enable, _ := engine.showLog[tag]; !enable {
		return
	}
	if fn, ok := LogTagProcessor[tag]; ok {
		format, contents = fn(tag, format, contents)
		if contents == nil {
			return
		}
	}
	engine.LogDebugf(format, contents...)
}

// logging warn
func (engine *Engine) TagLogWarn(tag string, contents ...interface{}) {
	if enable, _ := engine.showLog[tag]; !enable {
		return
	}
	if fn, ok := LogTagProcessor[tag]; ok {
		_, contents = fn(tag, "", contents)
		if contents == nil {
			return
		}
	}
	engine.LogWarn(contents...)
}

func (engine *Engine) TagLogWarnf(tag string, format string, contents ...interface{}) {
	if enable, _ := engine.showLog[tag]; !enable {
		return
	}
	if fn, ok := LogTagProcessor[tag]; ok {
		format, contents = fn(tag, format, contents)
		if contents == nil {
			return
		}
	}
	engine.LogWarnf(format, contents...)
}
