/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc     终端打印日志
* @name
* @author   huajie <huajie@baidu.com>
* @time     2018-08-24
*/

package writer

import (
    "io"
    "os"
    "fmt"
    "go-logger/x1"
)

// terminal print color
var (
    green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
    white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
    yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
    red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
    blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
    magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
    cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
    reset   = string([]byte{27, 91, 48, 109})
)

func ConsoleGreen(str string) string {
    return ConsoleColor(green, str)
}

func ConsoleWhite(str string) string {
    return ConsoleColor(white, str)
}

func ConsoleYellow(str string) string {
    return ConsoleColor(yellow, str)
}

func ConsoleRed(str string) string {
    return ConsoleColor(red, str)
}

func ConsoleBlue(str string) string {
    return ConsoleColor(blue, str)
}

func ConsoleMagenta(str string) string {
    return ConsoleColor(magenta, str)
}

func ConsoleCyan(str string) string {
    return ConsoleColor(cyan, str)
}

func ConsoleReset(str string) string {
    return ConsoleColor(reset, str)
}

func ConsoleColor(colorLeft, str string) string {
    return colorLeft + str + reset
}

// terminal print color
var stdout io.Writer = os.Stdout

// type LogWriterConsole LogWriterParent

/**
 * @desc   命令行输出日志
 */
type LogWriterConsole struct {
    /**
     * @desc    通道，传输LogRecord
     */
    rec chan *LogRecord
    /**
     *  @desc writer
     */
    LogWriter

}

/**
 * @desc    初始化命令行输出
 */
func NewLogWriterConsole() *LogWriterConsole{
    lwc := &LogWriterConsole{
        rec: make(chan *LogRecord),
    }
    var timeStr     string
    var timeStrAt   int64

    go func(){
        // for recItem := range lwc.rec {
        //     if at := recItem.Created.UnixNano() / 1e9; at != timeStrAt {
        //         timeStr = recItem.Created.Format("2006-01-02 15:04:05")
        //     }
        //
        //     // 是否需要输出文件名
        //     var filename string
        //     if recItem.Name != "" {
        //         filename = ConsoleBlue(recItem.Name) + ": "
        //     } else {
        //         filename = ""
        //     }
        //
        //     level := recItem.LevelStr
        //     if x1.LevelDanger(level) {
        //         level = ConsoleRed(level)
        //     }
        //
        //     level  += ": "
        //     source := " " + recItem.Source + " "
        //
        //     fmt.Fprint(stdout, filename, level, timeStr, source, recItem.Message, "\n")
        // }


        for {
            select {
            case recItem, ok := <-lwc.rec:
                if !ok {
                    return
                }
                if at := recItem.Created.UnixNano() / 1e9; at != timeStrAt {
                    timeStr = recItem.Created.Format("2006-01-02 15:04:05")
                }

                // 是否需要输出文件名
                var filename string
                if recItem.Name != "" {
                    filename = ConsoleBlue(recItem.Name) + ": "
                } else {
                    filename = ""
                }

                level := recItem.LevelStr
                if x1.LevelDanger(level) {
                    level = ConsoleRed(level)
                }

                level  += ": "
                source := " " + recItem.Source + " "

                fmt.Fprint(stdout, filename, level, timeStr, source, recItem.Message, "\n")
            }
        }
    }()

    return lwc
}

/**
 * @desc    执行终端输出
 */
func (lwc *LogWriterConsole) run(out io.Writer) {
    var timeStr     string
    var timeStrAt   int64

    for recItem := range lwc.rec {
        if at := recItem.Created.UnixNano() /1e9; at != timeStrAt {
            timeStr, timeStrAt = recItem.Created.Format("2018-08-24 15:30:00"), at
        }

        // 是否需要输出文件名
        var filename string
        if recItem.Name != "" {
            filename = ConsoleBlue(recItem.Name) + ": "
        } else {
            filename = ""
        }

        level := recItem.LevelStr
        if x1.LevelDanger(level) {
            level = ConsoleRed(level)
        }

        level  += ": "
        source := " " + recItem.Source + " "

        fmt.Fprint(out, filename, level, timeStr, source, recItem.Message, "\n")
    }
}

/**
 * @desc 实现接口写
 */
func (lwc *LogWriterConsole) LogWrite(rec *LogRecord) {
    if !x1.LogWithBlocking {
        if len(lwc.rec) >= x1.LogBufferLength {
            return
        }
    }

    lwc.rec <- rec
}

/**
 * desc 关闭通道
 */
func (lwc *LogWriterConsole) LogClose() {
    close(lwc.rec)
}


/**
 * @desc 实现初始化，并执行日志输出
 */
func RunLogWriterConsole(){
    lwc := NewLogWriterConsole()
    go lwc.run(stdout)
}



