# go-logger
## goland实现的log模块，可以在大型工程应用中，作为模块化，直接引入使用
### 实现功能：
1、log可配置化<br>
2、输出可根据错误级别自动写入${appname}.log 或者${appname}.log.wf <br>
3、输出支持输出到命令行终端 <br>
4、日志文件可以根据log配置中切割周期、切割数量、保留份数来维护服务器上日志<br>

### 文件目录：
#### 一、基本属性：
go-logger/x1/<br>
|文件名|说明|
|:---|:---|
|LogId|日志ID|
|LogLevel|日志级别|
|LogFormat|日志格式化方式|
|LogConfig|日志配置，维护configPool: 模块名->LogConfig|

| 日期      |    状态码  | 状态码总量 |(非)全量PV | 对应占比|
| --------: | --------:| ----:| --------: | --------:|
| 6月19日    | 499 |   7135 | 53585140 | 0.00013 |
| 6月19日    | 200 |   53577983 | 53585140 | 0.99987 |
| 6月20日    | 499 |   14331 | 94110198 | 0.00015 |
| 6月20日    | 200 |   94095746 | 94110198 | 0.99985 |
| 6月21日    | 499 |   16922 | 94803909 | 0.00018 |
| 6月21日    | 200 |   94786467 | 94803909 | 0.99982 |
| 6月22日    | 499 |   6557 | 41334571 | 0.00016 |
| 6月22日    | 200 |   41328005 | 41334571 | 0.99984 |

#### 二、基本操作：
LoggerPool      维护MAP: appname-->Logger<br>
Logger          维护整个logger对象，通过不同的标识映射不同的LoggerWriter<br>
LoggerWriter    维护日志的写，具体实现可以文本写，控制台输出等<br>
LoggerRecord    维护日志每行记录数据<br>

#### 三、例子
```
logObj := go_logger.LogInit("demo1")<br>
logObj.Error("errmsg","xxxdsfd")<br>
````
1、demo1是go-logger/config/log_conf.toml中配置，来初始Logger，无则通过LogConfig完成实例化<br>
2、调用关系logObj.Error("xxx") --> Logf() --> logObj.Logger.Write(rec LoggerRecord)<br>
