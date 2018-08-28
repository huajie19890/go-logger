/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc     日志LOG_LEVEL
* @name     LogLevel
* @author   huajie <huajie@baidu.com>
* @time     2018-08-07
*/


package x1

import (
    "strings"
)

/**
 * @desc LogLevel结构体定义
 */
type LogLevel struct {
}

/**
 * @desc log level const int 定义
 */
const (
    LOG_DEBUG  = iota // 0
    LOG_TRACE
    LOG_INFO
    LOG_NOTICE
    LOG_WARNING
    LOG_ERROR
    LOG_FATAL
)

/**
 * @desc log 聚合成不同输出文件
 * @desc normal   .log
 * @desc careful  .log.wf
 * @desc stdconsole 终端输出
 */
const (
    TXT_NORMAL  =   "normal"
    TXT_CAREFUL =   "careful"
    STD_CONSOLE =   "stdconsole"
)

/**
 * @desc log level string 枚举
 */
var levelString = []string{
    "DEBUG",
    "TRACE",
    "INFO",
    "NOTICE",
    "WARNING",
    "ERROR",
    "FATAL",
}

/**
* @desc log level 类型
* @desc 日志类型对应的日志级别
*/
var levelSwitch = map[int]string {
    LOG_DEBUG   : "normal",
    LOG_TRACE   : "normal",
    LOG_INFO    : "normal",
    LOG_NOTICE  : "normal",
    LOG_WARNING : "careful",
    LOG_ERROR   : "careful",
    LOG_FATAL   : "careful",
}

/**
 * @desc    字符串level字面量转换成常量字面量
 * @name    LevelToInt
 * @param   string  levelStr
 * @return  int     levelInt
 */
func LevelToInt( levelStr string ) int {
    var levelInt int
    levelStrUp := strings.ToUpper(levelStr)

    switch levelStrUp {
    case "DEBUG":
        levelInt = LOG_DEBUG
    case "TRACE":
        levelInt = LOG_TRACE
    case "INFO":
        levelInt = LOG_INFO
    case "NOTICE":
        levelInt = LOG_NOTICE
    case "WARNING":
        levelInt = LOG_WARNING
    case "ERROR":
        levelInt = LOG_ERROR
    case "FATAL":
        levelInt = LOG_FATAL
    default:
        levelInt = LOG_INFO // 2 常规日志级别
    }
    return levelInt
}

/**
 * @desc    日志的级别获取对应的日志文件类型
 * @name    LevelSwitch
 * @param   string  levelStr
 * @return  string
 */
func LevelSwitch(levelStr string) string {
    levelInt := LevelToInt(levelStr)
    levelSwitchType, _ := levelSwitch[levelInt]
    return levelSwitchType
}

/**
 * @desc    确定日志的级别
 */
func LevelDanger(levelStr string) bool {
    levelInt := LevelToInt(levelStr)
    if levelInt >= 5 {
        return true
    }
    return false
}














