/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc     日志LOG_ID
* @name     LogId
* @author   huajie <huajie@baidu.com>
* @time     2018-08-07
*/

package x1

import (
    "net/http"
    "strings"
    "strconv"
    "time"
)

/**
 * LogId结构体定义
 */
type LogId struct {
}


/**
 * @desc    生成logid
 * @desc    优先级表单中logid > header中logid > 随机生成
 * @name    CreateLogId
 * @param   request *http.Request
 * @return  ...     uint64
 */
func CreateLogId(request *http.Request) uint64 {

    var logidstr string
    var logidint uint64
    var err      error
    form := request.URL.Query()

    // 表单中的logid
    logidstr = strings.TrimSpace(form.Get("logid"))
    if logidstr != "" {
        if logidint, err = strconv.ParseUint(logidstr, 10, 64); err == nil && logidint > 0 {
            return logidint
        }
    }

    // header中的logid API上下游透传
    logidstr = request.Header.Get("logid")
    if logidstr != "" {
        if logidint, err = strconv.ParseUint(logidstr,10, 64); err == nil && logidint > 0 {
            return logidint
        }
    }

    // 随机生成
    usec    := uint64(time.Now().UnixNano())
    logidint = usec&0x7FFFFFFF | 0x80000000
    return logidint
}












