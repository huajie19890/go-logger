/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc
* @name     Logger
* @author   huajie <huajie@baidu.com>
* @time     2018-08-07
*/


package go_logger

import (
    "time"
    "runtime"
    fmt "fmt"
    "strings"
    "go-logger/utils"
    "go-logger/writer"
    "go-logger/x1"
    "os"
)

/**
 * Logger struct 定义
 */
type Logger struct {
    /**
     * @desc    logger对象的名称 例如demo\dblog
     */
    LoggerName  string
    /**
     * @desc    loggerconf
     */
    LoggerConf  *x1.LogConfig
    /**
     * @desc    log包含的Loggerchooser，不同级别的日志选择写入不同的日志文件
     */
    Chooser     map[string]*LoggerChooser
}


/**
 * @desc 初始化Logger
 */
func NewLogger(conf *x1.LogConfig) *Logger{

    // 避免重复创建
    if loggerObj, ok := LoggerPool[conf.LogConfName]; ok {
        if len(loggerObj.Chooser) > 0 {
            return loggerObj
        }
        delete (LoggerPool, conf.LogConfName)
        loggerObj.Close()
    }

    loggerObj := &Logger{
        LoggerName: conf.LogConfName,
        LoggerConf: conf,
        Chooser:    make(map[string]*LoggerChooser),
    }

    // 初始化日志路径
    if err := initLogDir(conf); err != nil {
        panic(err)
    }
    fmt.Println("logobj:(name,conf,chooser)" , loggerObj, "logobj地址:", &loggerObj)
    // 实例化 TXT常规日志
    lcTxtNorObj := NewLoggerChooserTxt(conf, x1.TXT_NORMAL)
    loggerObj.Chooser[x1.TXT_NORMAL] = lcTxtNorObj

    // 实例化 TXT告警日志
    lcTxtCareObj:= NewLoggerChooserTxt(conf, x1.TXT_CAREFUL)
    loggerObj.Chooser[x1.TXT_CAREFUL] = lcTxtCareObj
    // 实例化 终端输出日志
    if conf.Stdout {
        lcConsoleObj := NewLoggerChooserConsole(conf, x1.STD_CONSOLE)
        loggerObj.Chooser[x1.STD_CONSOLE] = lcConsoleObj
    }
    // 维护全局MAP
    LoggerPool[conf.LogConfName] = loggerObj
    // 返回logger对象
    return loggerObj
}



/**
 * @desc 支持 DBUEG 打印日志
 */
func (logger *Logger) Debug(arg0 interface{}, args... interface{}) {
    logger.gather("DEBUG", arg0, args...)
}

/**
 * @desc 支持 TRACE 打印日志
 */
func (logger *Logger) Trace(arg0 interface{}, args... interface{}) {
    logger.gather("TRACE", arg0, args...)
}

/**
 * @desc 支持 INFO 打印日志
 */
func (logger *Logger) Info(arg0 interface{}, args... interface{}) {
    logger.gather("INFO", arg0, args...)
}

/**
 * @desc 支持 NOTICE 打印日志
 */
func (logger *Logger) Notice(arg0 interface{}, args... interface{}) {
    logger.gather("NOTICE", arg0, args...)
}

/**
 * @desc 支持 ERROR 打印日志
 */
func (logger *Logger) Error(arg0 interface{}, args... interface{}) {
    logger.gather("ERROR", arg0, args...)
}

/**
 * @desc 支持 FATAL 打印日志
 */
func (logger *Logger) Fatal(arg0 interface{}, args... interface{}) {
    logger.gather("FATAL", arg0, args...)
}


/**
 * @desc    各种日志打印集中处理
 * @name    gather
 * @param   interface args...
 * @return  nil
 */
func (logger *Logger) gather(levelStr string, arg0 interface{}, args... interface{}) {
    switch first := arg0.(type) {
    case string:
        logger.logf(levelStr, first, args...)
    case []byte:
        logger.logf(levelStr, string(first))
    case map[string]string:
        logger.logf(levelStr, utils.LogNsdFormat(first, logger.Separator()), args...)
    case [][2]string:
        logger.logf(levelStr, utils.LogSdFormat(first, logger.Separator()), args...)
    case func() string:
        logger.logf(levelStr, first())
    default:
        logger.logf(levelStr, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
    }
}


/**
 * @desc    获取日志配置中分隔符
 * @name    Separator
 * @param   nil
 * @return  nil
 */
func (logger *Logger) Separator() string {
    return logger.LoggerConf.Separator
}


/**
 * @desc    获取日志配置中caller
 * @name    SkipCaller
 * @param   nil
 * @return  nil
 */
func (logger *Logger) SkipCaller() int {
    return logger.LoggerConf.SkipCaller
}


/**
 * @desc    日志格式化,给logrecod赋值
 * @name    logf
 * @param   string  levelStr
 * @param   string  format
 * @param   interface{} args...
 * @return  nil
 */
func (logger *Logger) logf(levelStr string, format string, args... interface{}) {


    src := ""
    skipTimes := 3
    skipTimes  =  logger.SkipCaller() + skipTimes
    if pc, _, _, ok := runtime.Caller(skipTimes); ok {
        filename, lineno := runtime.FuncForPC(pc).FileLine(pc)
        caller := map[string]string{"caller": fmt.Sprintf("%s:%d", filename, lineno)}
        src = utils.LogNsdFormat(caller, logger.Separator())
    }

    msg := format
    if len(args) > 0 {
        msg = fmt.Sprintf(format + logger.Separator() + "%v", args...)
    }


    rec := &writer.LogRecord{
        Name:       "",
        LevelStr:   levelStr,
        LevelInt:   x1.LevelToInt(levelStr),
        ProcessId:  "",
        Created:    time.Now(),
        Source:     src,
        Message:    msg,
    }

    levelSwitch := x1.LevelSwitch(levelStr)
    logger.Chooser[levelSwitch].LogWriterObj.LogWrite(rec)

    if logger.LoggerConf.Stdout {
        logger.Chooser[x1.STD_CONSOLE].LogWriterObj.LogWrite(rec)
    }
}


/**
 * @desc    清除Logger一切，包括chooser, writer
 */
func (logger *Logger) Close() {
    for chooserName, chooserObj := range logger.Chooser {
        chooserObj.Close()
        delete(logger.Chooser, chooserName)
    }
}



/**
 * @desc 初始化日志路径
 */
func initLogDir(conf *x1.LogConfig) error {
    path := conf.LogPath
    if path == "" {
        path = "./"
    }
    _, statErr := os.Stat(path)
    if os.IsNotExist(statErr) {
        ctErr := os.MkdirAll(path, 0777)
        if ctErr != nil {
            return fmt.Errorf("log conf err: create log dir '%s' error: %s", path, ctErr)
        }
    }
    return nil
}
