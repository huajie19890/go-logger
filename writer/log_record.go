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
    "time"
)

type LogRecord struct {
    /**
     * @desc    The log file which will be writen into
     */
    Name        string
    /**
     * @desc    The log LevelStr
     */
    LevelStr    string
    /**
     * @desc    The log LevelInt
     */
    LevelInt    int
    /**
     * @desc    ProcessId
     */
    ProcessId   string
    /**
     * @desc    The time at which the log message was created (nanoseconds)
     */
    Created     time.Time
    /**
    * @desc    The message source
    */
    Source      string
    /**
    * @desc     The log message
    */
    Message     string
}








