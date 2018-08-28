/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc     接口定义
* @name     LogWriter
* @author   huajie <huajie@baidu.com>
* @time     2018-08-21
*/


package writer

type LogWriter interface {
    /**
     * @desc 记录日志
     */
    LogWrite(rec *LogRecord)

    /**
     * @desc 关闭日志writer
     */
    LogClose()
}



