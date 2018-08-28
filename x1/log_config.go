/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc     日志配置处理
* @name     LogConfig
* @author   huajie <huajie@baidu.com>
* @time     2018-08-07
*/


package x1

import (
    "strings"
    "log"
)

/**
 * @desc    LogConfig对象申明
 * @name    LogConfig
 */
type LogConfig struct {
    /**
     * @desc 调用日志模块的app名字、或者是业务名称作为LogConfName标识
     */
    LogConfName string
    /**
     * @desc 日志打印路径
     */
    LogPath    string
    /**
     * @desc 日志切割周期\速率
     */
    RotateRate  string
    /**
     * @desc 日志轮转切割次数
     */
    RotateCount int
    /**
     * @desc    日志类型STRING
     */
    LevelStr    string
    /**
     * @desc    日志类型INT
     */
    LevelInt    int
    /**
     * @desc    日志格式化
     */
    // Format      string
    /**
     * @desc    分割符号 logid[12xxxx] logid=12xxxx
     */
    Separator   string
    /**
     * @desc    是否同时console打印
     */
    Stdout      bool
    /**
     * @desc    栈信息打印跳过次数
    */
    SkipCaller  int
}

/**
 * @desc 全局的LOG变量-配置
 */
 var (
    // LogBufferLength 指定每一个message的长度
    LogBufferLength = 4096
    // 当buffer满载时是否锁住，若通道打满则会锁住进程
    LogWithBlocking = true
    // 默认日志格式
    LogFormat = FORMAT_DEFAULT_SEC
    // 默认进程ID
    LogProcessId = "0"
    // 当打印二进制文件时是否输出打印源码行位置 此配置目前无效
    EnableSrcForBinLog = true
)


/**
 * @desc 默认的配置
 */
var LogConfigDefault = LogConfig{
    LogConfName:    "unkown",
    LogPath:        "./",
    RotateRate:     "h",
    RotateCount:    48,
    LevelStr:       "INFO",
    LevelInt:       LOG_INFO,
    Separator:      "=",
    Stdout:         true,
    SkipCaller:     0,
}

/**
 * @desc    设置日志配置
 * @name    初始化NewLogConfig
 * @param   LogConfig   myConf
 * @return  *LogConfig
 */
func NewLogConfig(myConf LogConfig) *LogConfig {
    // 初始化配置名字
    if myConf.LogConfName == "" {
        myConf.LogConfName = LogConfigDefault.LogConfName
    }

    // 初始化日志路径
    if myConf.LogPath == "" {
        myConf.LogPath = LogConfigDefault.LogPath
    }

    // 初始化日志切割速率
    if myConf.RotateRate == "" {
        myConf.RotateRate = LogConfigDefault.RotateRate
    }
    supportRate := map[string]string{"M":"1", "H":"2", "D":"3","":"4",}
    if _, exists := supportRate[strings.ToUpper(myConf.RotateRate)]; !exists {
        panic("日志切割不在可支持的列表中！请修改配置：D或H或M")
    }

    // 初始化日志保留周期
    if myConf.RotateCount < 0 {
        myConf.RotateCount = LogConfigDefault.RotateCount
    }

    // 初始化日志级别
    if myConf.LevelStr == "" {
        myConf.LevelStr = LogConfigDefault.LevelStr
    }

    // 初始化日志级别int
    if myConf.LevelInt < 0 {
        myConf.LevelInt = LevelToInt(myConf.LevelStr)
    }

    // 初始化日志msg中 key value 分割标记
    if myConf.Separator == "" {
        myConf.Separator = LogConfigDefault.Separator
    }

    // 初始化是否支持终端输出日志
    if myConf.Stdout == true {
        myConf.Stdout = true
    } else {
        myConf.Stdout = false
    }

    // 初始化追踪栈的层数
    if myConf.SkipCaller < 0 {
        myConf.SkipCaller = LogConfigDefault.RotateCount
    }
    return &myConf
}





// ====================  LogConfigs ==================== //

/**
 * @desc    存放log config
 */
type LogConfigs struct {
    /**
     * @desc 全局存放LogConfig映射
     */
    LogConfigPool  map[string]*LogConfig
}



/**
 * @desc init LogConfigs
 */
func NewLogConfigs(myConfs []LogConfig) *LogConfigs{
    logconfigs := &LogConfigs{
        LogConfigPool: make(map[string]*LogConfig),
    }
    logconfigs.SetLogConfigs(myConfs)
    return logconfigs
}


/**
 * @desc    设置多个日志配置
 * @name    SetLogConfigs
 * @param   []LogConfig   myConfs
 * @return  LogConfigPool
 */
func (logConfigs *LogConfigs)SetLogConfigs(myConfs []LogConfig) map[string]*LogConfig {
    if len(myConfs) > 0 {
        for _, eachConf := range myConfs {
            if eachConf.LogConfName == "" {
                log.Fatal("log_conf配置有误，LogConfName为空！")
                continue
            }

            if _, isExist := logConfigs.LogConfigPool[eachConf.LogConfName]; isExist {
                delete(logConfigs.LogConfigPool, eachConf.LogConfName)
            }
            latestConf := NewLogConfig(eachConf)
            logConfigs.LogConfigPool[eachConf.LogConfName] = latestConf
        }
    } else {  // 配置文件为空情况, 默认值赋值
        logConfigs.LogConfigPool[LogConfigDefault.LogConfName] = &LogConfigDefault
    }
    return logConfigs.LogConfigPool
}

/**
 * @desc    获取多个日志配置
 * @name    GetLogConfigPool
 * @param   []LogConfig myConfs
 * @return  map[string]*LogConfig   LogConfigPool
 */
func (logConfigs *LogConfigs)GetLogConfigPool(myConfs []LogConfig) map[string]*LogConfig {
    logConfigs.LogConfigPool = logConfigs.SetLogConfigs(myConfs)
    return logConfigs.LogConfigPool
}








