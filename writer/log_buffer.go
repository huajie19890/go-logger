/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc     对每一条记录管理
* @name     LogRecord
* @author   huajie <huajie@baidu.com>
* @time     2018-08-21
*/


package writer

import (
    "sync"
    "bytes"
    "fmt"
)

type LogBuffer struct {
    /**
     * @desc 最近一次更新记录的时间
     */
    LastUpdateSeconds    int64
    /**
     * @desc 简短时间、日期
     */
    shortTime, shortDate string
    /**
     * @desc 长时间、日期
     */
    longTime, longDate   string

}

/**
 * @desc    获取LogBuffer对象
 */
var logBuffer = &LogBuffer{}
/**
 * @desc    格式化用到的锁
 */
var formatMutex sync.Mutex
/**
 * @desc    临时对象池
 */
var bufPool     sync.Pool


/**
 * @desc 创建一个4KB的缓存
 */
func newBuf() *bytes.Buffer {
    if v := bufPool.Get(); v != nil {
        return v.(*bytes.Buffer)
    }
    return bytes.NewBuffer(make([]byte, 0, 4096))
}

/**
 * @desc 放入数据
 */
func putBuf(byteBuffer *bytes.Buffer) {
    // 重置缓冲区为空
    byteBuffer.Reset()
    bufPool.Put(byteBuffer)
}

/**
 * @desc    format LogRecord
 * @name    LogRecordFormat
 * @param   string      format
 * @param   *LogRecord  rec
 * @return  string
 *
 * Known format codes:
 * %T - Time (15:04:05 MST)
 * %t - Time (15:04)
 * %D - Date (2006/01/02)
 * %d - Date (01/02/06)
 * %L - Level (DEBG, TRAC, INFO, NOTICE, WARNING, ERROR, FATAL)
 * %P - Pid of process
 * %S - Source
 * %M - Message
 * Ignores unknown formats
 * Recommended: "[%D %T] [%L] (%S) %M"
 *
 */
func LogRecordFormat(format string, rec *LogRecord) string {

    if rec == nil {
        return "<nil>"
    }
    if len(format) == 0 {
        return ""
    }

    out := newBuf()
    defer putBuf(out)

    secs := rec.Created.UnixNano() / 1e9

    formatMutex.Lock()
    cache := *logBuffer
    formatMutex.Unlock()
    if cache.LastUpdateSeconds != secs {
        month, day, year := rec.Created.Month(), rec.Created.Day(), rec.Created.Year()
        hour, minute, second := rec.Created.Hour(), rec.Created.Minute(), rec.Created.Second()
        zone, _ := rec.Created.Zone()
        updated := &LogBuffer{
            LastUpdateSeconds: secs,
            shortTime:         fmt.Sprintf("%02d:%02d", hour, minute),
            shortDate:         fmt.Sprintf("%02d/%02d/%02d", month, day, year%100),
            longTime:          fmt.Sprintf("%02d:%02d:%02d %s", hour, minute, second, zone),
            longDate:          fmt.Sprintf("%04d/%02d/%02d", year, month, day),
        }
        formatMutex.Lock()
        cache = *updated
        logBuffer = updated
        formatMutex.Unlock()
    }

    // Split the string into pieces by % signs
    pieces := bytes.Split([]byte(format), []byte{'%'})

    // Iterate over the pieces, replacing known formats
    for i, piece := range pieces {
        if i > 0 && len(piece) > 0 {
            switch piece[0] {
            case 'T':
                out.WriteString(cache.longTime)
            case 't':
                out.WriteString(cache.shortTime)
            case 'D':
                out.WriteString(cache.longDate)
            case 'd':
                out.WriteString(cache.shortDate)
            case 'L':
                out.WriteString(rec.LevelStr)
            case 'P':
                out.WriteString(rec.ProcessId)
            case 'S':
                out.WriteString(rec.Source)
            case 'M':
                out.WriteString(rec.Message)
            }
            if len(piece) > 1 {
                out.Write(piece[1:])
            }
        } else if len(piece) > 0 {
            out.Write(piece)
        }
    }
    out.WriteByte('\n')

    return out.String()
}