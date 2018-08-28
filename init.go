/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc     入口函数
* @name     init
* @author   huajie <huajie@baidu.com>
* @time     2018-08-07
*/


package go_logger

import (
    "fmt"
    "github.com/BurntSushi/toml"

    "go-logger/x1"
    "sync"
)

// @todo 1.处理日志配置 logInit("demo")

// @todo 2.激活 logger NewLogger(conf)

// @todo 3.写入 logger.Info("msg")  // logger.Fatal("msg")

// logObj := LogInit("demo")
// logObj.Info("hello")
// logObj.Fatal("world")

var once        sync.Once
var myLogger    *Logger

func LogInit(loggerName string) *Logger{

    once.Do(func() {
        myLogger = logInit(loggerName)
    })
    return myLogger
}


func logInit(loggerName string) *Logger{

    var myConf map[string][]x1.LogConfig
    _, err := toml.DecodeFile("./../config/log_conf.toml", &myConf)
    fmt.Println("Log init start------loadConf:", myConf, "error:", err)
    myConfs, _ := myConf["logconfig"]

    logConfigs := x1.NewLogConfigs(myConfs)
    if  loggerConfig, exists := logConfigs.LogConfigPool[loggerName]; exists{
        logger := NewLogger(loggerConfig)
        return logger
    } else {
       loggerConfig := x1.LogConfigDefault
       logger := NewLogger(&loggerConfig)
       return logger
    }
}


