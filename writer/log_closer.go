/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc     处理日志关闭等
* @name     LogCloser
* @author   huajie <huajie@baidu.com>
* @time     2018-08-22
*/

package writer

type LogCloser struct {
    /**
     * @desc    chan判断是否结束
     */
    isEnd chan bool
}


/**
 * @desc   初始化
 */
// func NewLogCloser() *LogCloser {
//     isEnd := make(chan bool)
//     return &LogCloser{isEnd:isEnd}
// }

func (logCloser *LogCloser) LogCloserInit() {
    logCloser.isEnd = make(chan bool)
}

/**
 * @desc    通知该结束了
 */
func (logCloser *LogCloser) EndNotify(rec *LogRecord) bool{
    if rec == nil && logCloser != nil {
        logCloser.isEnd <- true
        return true
    }
    return false
}

/**
 * @desc    等待进行结束
 */
func (logCloser *LogCloser) EndWait(rec chan *LogRecord) {
    rec <- nil
    if logCloser.isEnd != nil {
        <-logCloser.isEnd
    }
}
