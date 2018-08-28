/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc     日志格式化
* @name     LogFormat
* @author   huajie <huajie@baidu.com>
* @time     2018-08-13
*/


package x1


const (
    FORMAT_DEFAULT          = "[%D %T] [%L] (%S) %M"
    FORMAT_DEFAULT_WITH_PID = "[%D %T] [%L] [%P] (%S) %M"
    FORMAT_SHORT            = "[%t %d] [%L] %M"
    FORMAT_ABBREV           = "[%L] %M"
    FORMAT_DEFAULT_SEC      = "%L: %D %t %S %M"
)
