/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc     写文件定义
* @name     LogWriterTxt
* @author   huajie <huajie@baidu.com>
* @time     2018-08-22
*/


package writer

import (
    "os"
    "regexp"
    "strings"
    "time"
    "fmt"
    "path/filepath"
    "io/ioutil"
    "sort"
    "go-logger/x1"
    "go-logger/utils"
)

const (
    /* number of seconds in a day */
    MIDNIGHT = 24 * 60 * 60
    /* number of seconds in a hour */
    NEXTHOUR = 60 * 60
)



type LogWriterTxt struct {

    /**
     * @desc    通道，传输LogRecord
     */
    rec         chan *LogRecord
    /**
     * @desc    日志存放基础路径
     */
    basename    string
    /**
     * @desc    日志名称
     */
    filename    string
    /**
     * @desc    日志全路径
     */
    baseFilename string
    /**
     * @desc    日志文件句柄
     */
    file         *os.File
    /**
     * @desc    日志格式化要求
     */
    format      string
    /**
     * @desc    日志文件后缀
     */
    suffix      string
    /**
     * @desc    'D', 'H', 'M', "MIDNIGHT", "NEXTHOUR" 轮询周期
     */
    when        string
    /**
     * @desc    备份日志文件数量
     */
    backupCount int
    /**
     * @desc    间隔
     */
    interval   int64
    /**
     * @desc    旧日志清理
     */
    fileFilter *regexp.Regexp
    /**
     * @desc    time.Unix()
     */
    rolloverAt int64
    /**
     * @desc    fixed caller stack skip times
     */
    callerSkip int
    /**
     * @desc    日志内细节分割标记
     */
    separator string
    /**
     * @desc 处理close
     */
    LogCloser
    /**
     *  @desc writer
     */
    LogWriter
}

/**
 * @desc    初始化
 */
func NewLogWriterTxt(fname string, config *x1.LogConfig) *LogWriterTxt {

    lwt := &LogWriterTxt{
       // rec:            make(chan *LogRecord, x1.LogBufferLength),
        rec:            make(chan *LogRecord),
        basename:       filepath.Base(fname),
        filename:       fname,
        format:         x1.LogFormat,
        when:           strings.ToUpper(config.RotateRate),
        backupCount:    config.RotateCount,
        callerSkip:     config.SkipCaller,
        separator:      config.Separator,
    }

    lwt.LogCloserInit()

    if path, err := filepath.Abs(fname); err != nil {
        fmt.Fprintf(os.Stderr, "NewLogWriterTxt(%q): %s\n", lwt.filename, err)
        return nil
    } else {
        lwt.baseFilename = path
    }

    // file prepare
    lwt.prepare()

    // open the file for the first time, then routate it if necessary
    if err := lwt.rolloverInit(); err != nil {
        fmt.Fprintf(os.Stderr, "NewLogWriterTxt(%q): %s\n", lwt.filename, err)
        return nil
    }

    // start  goroutine
    go func() {
        defer func() {
            if lwt.file != nil {
                lwt.file.Close()
            }
        }()

        for {
            select {
            case rec, ok := <-lwt.rec:
                if !ok {
                    return
                }

                if lwt.EndNotify(rec) {
                    return
                }

                if lwt.rolloverNeeded() {
                    if err := lwt.rolloverInit(); err != nil {
                        fmt.Fprintf(os.Stderr, "NewTimeFileLogWriter(%q): %s\n", lwt.filename, err)
                        return
                    }
                }

                // Perform the write
                var err error
                // fmt.Println(5, LogRecordFormat(lwt.format, rec))
                _, err = fmt.Fprint(lwt.file, LogRecordFormat(lwt.format, rec))
                if err != nil {
                    fmt.Fprintf(os.Stderr, "NewTimeFileLogWriter(%q): %s\n", lwt.filename, err)
                    return
                }

            }
        }
    }()

    return lwt
}


// ================   方法集  ==================== //
/**
 * @desc    实现写日志方法
 * @name    LogWrite
 * @param   *LogRecord  rec
 * @return  nil
 */
func (lwt *LogWriterTxt) LogWrite(rec *LogRecord) {
    if !x1.LogWithBlocking {
        if len(lwt.rec) >= x1.LogBufferLength {
            return
        }
    }
    lwt.rec <- rec
}

/**
 * @desc    关闭整个writer包括chan
 * @name    LogWrite
 * @param   *LogRecord  rec
 * @return  nil
 */
func (lwt *LogWriterTxt) LogClose() {
    lwt.EndWait(lwt.rec)
    close(lwt.rec)
}

/**
 * @desc    通过切割周期获取对应的日志落地文件
 * @name    prepare
 * @param   nil
 * @return  nil
 */
func (lwt *LogWriterTxt) prepare() {
    var regRule string

    switch lwt.when {
    case "M":  // 分钟级别
        lwt.interval = 60
        lwt.suffix = "%Y%m%d%H%M"
        regRule = `^\d{4}\d{2}\d{2}\d{2}\d{2}$`
    case "H": // 小时级别
        lwt.interval = 60 * 60
        lwt.suffix = "%Y%m%d%H"
        regRule = `^\d{4}\d{2}\d{2}\d{2}$`
    case "D": // 天级别
        lwt.interval = 60 * 60 * 24
        lwt.suffix = "%Y%m%d"
        regRule = `^\d{4}\d{2}\d{2}$`
    case "":
        lwt.interval = 0
        lwt.suffix = ""
    default: // default is "D"
        lwt.interval = 60 * 60 * 24
        lwt.suffix = "%Y%m%d"
        regRule = `^\d{4}\d{2}\d{2}$`
    }

    if lwt.interval != 0 {
        lwt.fileFilter = regexp.MustCompile(regRule)
        fInfo, err := os.Stat(lwt.filename)

        var t time.Time
        if err == nil {
            t = fInfo.ModTime() // 最后修改时间
        } else {
            t = time.Now()
        }
        lwt.rolloverAt = lwt.rolloverTime(t)
        return
    }
    lwt.rolloverAt = -1
}


/**
 * @desc    根据当前时间推算需要日志切割的时间戳
 * @name    rolloverTime
 * @param   time.Time   currTime
 * @return  int64
 */
func (lwt *LogWriterTxt) rolloverTime(currTime time.Time) int64 {
    var result int64

    if lwt.when == "D" {
        t := currTime.Local()
        /* r is the number of seconds left between now and midnight */
        r := MIDNIGHT - ((t.Hour()*60+t.Minute())*60 + t.Second())
        result = currTime.Unix() + int64(r)
    } else if  lwt.when == "H" {
        t := currTime.Local()
        /* r is the number of seconds left between now and the next hour */
        r := NEXTHOUR - (t.Minute()*60 + t.Second())
        result = currTime.Unix() + int64(r)
    } else if lwt.when == ""  {
        result = -1
    } else {
        result = currTime.Unix() + lwt.interval
    }
    return result
}


/**
 * @desc    通过当前时间判断日志是否需要切割，用于轮询机制中
 * @name    rolloverNeeded
 * @param   nil
 * @return  bool
 */
func (lwt *LogWriterTxt) rolloverNeeded() bool {
    if lwt.rolloverAt == -1 {
        return false
    }

    t := time.Now().Unix()
    if t >= lwt.rolloverAt {
        return true
    } else {
        return false
    }
}

/**
 * @desc    日志切割初始化
 * @name    rolloverInit
 * @param   nil
 * @return  error
 */
func (lwt *LogWriterTxt) rolloverInit() error {
    // Close any log file that may be open
    if lwt.file != nil {
        lwt.file.Close()
    }

    // check need rollover . then rename file to backup name
    if lwt.rolloverNeeded() {
        if err := lwt.rolloverBackup(); err != nil {
            return err
        }
    }

    // remove files, according to backupCount
    if lwt.backupCount > 0 {
        for _, fileName := range lwt.rolloverDelete() {
            os.Remove(fileName)
        }
    }

    // Open the log file
    fd, err := os.OpenFile(lwt.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil {
        return err
    }
    lwt.file = fd

    // adjust rolloverAt
    lwt.rolloverTimeAdjust()

    return nil
}


/**
 * @desc    日志切割-切割备份
 * @name    rolloverBackup
 * @param   nil
 * @return  error
 */
func (lwt *LogWriterTxt) rolloverBackup() error {
    _, err := os.Lstat(lwt.filename)
    if err == nil { // file exists
        // 推算出之切割时间之前的最近一个时间切割点
        t := time.Unix(lwt.rolloverAt-lwt.interval, 0).Local()
        fname := lwt.baseFilename + "." + utils.TimeFormat(lwt.suffix, t)

        // 如果重名文件，移除
        if _, err := os.Stat(fname); err == nil {
            err = os.Remove(fname)
            if err != nil {
                return fmt.Errorf("rolloverBackup: %s\n", err)
            }
        }

        // 重新命名
        err = os.Rename(lwt.baseFilename, fname)
        if err != nil {
            return fmt.Errorf("rolloverBackup: %s\n", err)
        }
    }
    return nil
}

/**
 * @desc    日志切割-获取需要删除的日志
 * @name    rolloverBackup
 * @param   nil
 * @return  error
 */
func (lwt *LogWriterTxt) rolloverDelete()  []string {
    dirName  := filepath.Dir(lwt.baseFilename)
    baseName := filepath.Base(lwt.baseFilename)
    result   := []string{}

    fileInfos, err := ioutil.ReadDir(dirName)
    if err != nil {
        fmt.Fprintf(os.Stderr, "rolloverDelete(%q): %s\n", lwt.filename, err)
        return result
    }

    prefix := baseName + "."
    plen   := len(prefix)

    for _, fileInfo := range fileInfos {
        fileName := fileInfo.Name()
        if len(fileName) >= plen && fileName[:plen] == prefix {
            suffix := fileName[plen:]
            // 若不需要分割时 fileFilter为空
            if lwt.fileFilter != nil && lwt.fileFilter.MatchString(suffix) {
                result = append(result, filepath.Join(dirName, fileName))
            }
        }
    }

    sort.Sort(sort.StringSlice(result))

    if len(result) < lwt.backupCount {
        result = result[0:0]
    } else {
        result = result[:len(result)-lwt.backupCount]
    }
    return result
}


/**
 * @desc    日志切割-修正未来切割时间
 * @name    rolloverBackup
 * @param   nil
 * @return  error
 */
func (lwt *LogWriterTxt) rolloverTimeAdjust() {
    if lwt.interval == 0 {
        return
    }

    currTime := time.Now()
    newRolloverAt := lwt.rolloverTime(currTime)

    for newRolloverAt <= currTime.Unix() {
        newRolloverAt = newRolloverAt + lwt.interval
    }

    lwt.rolloverAt = newRolloverAt
}


/**
 * @desc    获取缓存chann长度
 * @name    QueueLen
 * @param   nil
 * @return  int
 */
func (lwt *LogWriterTxt) QueueLen() int {
    return len(lwt.rec)
}

