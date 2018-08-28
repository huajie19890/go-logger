package main

import (
    "go-logger"

    "fmt"
)

func main() {

    logObj := go_logger.LogInit("demo1")


   //  var j int
   //  for  j=1; j<=800; j++ {
   //      str :=  fmt.Sprintf("index: %d", j)
   //      logObj.Info(str)
   //  }

    chanDone := make(chan bool, 1)

    go func(){
        var k int
        for k =1 ;  k<=7 ; k++ {
            str2 :=  fmt.Sprintf("index:%d", k)
            logObj.Info(str2)
        }
        chanDone <- true
    }()

    <- chanDone

    var i int
    for  i=1; i<=10; i++ {
        logObj.Error("w1","10")

    }
}