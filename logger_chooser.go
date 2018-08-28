/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc
* @name     LoggerChooser
* @author   huajie <huajie@baidu.com>
* @time     2018-08-22
*/


package go_logger

import (
    "fmt"
    "strings"
    "path/filepath"

    "go-logger/x1"
    "go-logger/writer"
)

type LoggerChooser struct {
    LogChooserName  string
    LogWriterObj    writer.LogWriter

}


/**
 * @desc    初始化TXT类型的对象 (分两个)
 * @name    NewLoggerChooserTxt
 * @param   *x1.LogConfig   conf
 * @param   string  levelSwitch
 * @return  *LoggerChooser
 */
func NewLoggerChooserTxt(conf *x1.LogConfig, levelSwitch string) *LoggerChooser{

    filename     := fileNameFullCreate(conf.LogConfName, conf.LogPath, levelSwitch)
    logWriterTxt := writer.NewLogWriterTxt(filename, conf)
    if logWriterTxt == nil {
        panic(fmt.Errorf("log conf err: in NewLoggerChooserTxtNoraml(%s)", filename))
    }

    return   &LoggerChooser{
        LogChooserName: levelSwitch,
        LogWriterObj:   logWriterTxt,
    }
}


/**
 * @desc    初始化STD类型的对象 (如果需要的话)
 * @name    NewLoggerChooserTxt
 * @param   *x1.LogConfig   conf
 * @param   string  levelSwitch
 * @return  *LoggerChooser
 */
func NewLoggerChooserConsole(conf *x1.LogConfig, levelSwitch string) *LoggerChooser{

    return &LoggerChooser{
        LogChooserName: levelSwitch,
        LogWriterObj:   writer.NewLogWriterConsole(),
    }
}


/**
 * @desc    write相关资源关闭
 * @name    Close
 * @param   nil
 * @return  nil
 */
func (loggerChooser *LoggerChooser) Close() {
    loggerChooser.LogWriterObj.LogClose()
}


//  ==================   函数定义 ===================  //
/**
 * @desc 生成文件全路径
 * loggerName :testapp
 * logDir: /home/work/demo/log
 * logLevelSwitch : careful
 * output : /home/work/demo/log/testapp.log.wf
 */
func fileNameFullCreate(loggerName string, logDir string, logLevelSwitch string) string {
    var suffix string
    switch logLevelSwitch {
    case x1.TXT_NORMAL:
        suffix = ".log"
    case x1.TXT_CAREFUL:
        suffix = ".log.wf"
    default:
        suffix = ".log"
    }
    strings.TrimSuffix(logDir, "/")
    return filepath.Join(logDir, loggerName + suffix)
}
