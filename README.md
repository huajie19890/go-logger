# go-logger
## goland实现的log模块，可以在大型工程应用中，作为模块化，直接引入使用
### 实现功能：
1、log可配置化<br>
2、输出可根据错误级别自动写入${appname}.log 或者${appname}.log.wf <br>
3、输出支持输出到命令行终端 <br>
4、日志文件可以根据log配置中切割周期、切割数量、保留份数来维护服务器上日志<br>

### 文件目录：
#### 一、基本属性：
go-logger/x1/

 文件名    |    说明  
 -------- |  --------
LogId    | 日志ID   
LogLevel    | 日志级别 
LogFormat    | 日志格式化方式
LogConfig    | 日志配置，维护configPool: 模块名->LogConfig 


#### 二、基本操作：

 文件名    |    说明  
 -------- |  --------
LoggerPool  |    维护MAP: appname-->Logger<br>
Logger      |    维护整个logger对象，通过不同的标识映射不同的LoggerWriter<br>
LoggerWriter |   维护日志的写，具体实现可以文本写，控制台输出等<br>
LogWriterTxt |   维护文本写，实现LoggerWriter接口
LogWriterConsole|维护终端写，实现LoggerWriter接口
LoggerRecord |   维护日志每行记录数据<br>

#### 三、例子
```
logObj := go_logger.LogInit("demo1")
logObj.Error("errmsg","xxxdsfd")
````
1、"demo1"是go-logger/config/log_conf.toml中配置，通过配置来初始Logger，无配置则通过默认的LogConfig完成实例化<br>
2、go_logger.LogInit("demo1")完成了日志的句柄的初始化，通过map维护在内存中
3、logObj.Error("xxx")调用关系为:logObj.Error("xxx") --> Logf() --> logger.Choose[levelSwitch].LogWriterObj.LogWrite(rec LoggerRecord)<br>
