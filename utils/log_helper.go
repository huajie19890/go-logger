/***************************************************************************
*                                                                          *
* Copyright (c) 2018 , Inc. All Rights Reserved                            *
*                                                                          *
****************************************************************************
*
* @desc
* @name     LogHelper
* @author   huajie <huajie@baidu.com>
* @time     2018-08-24
*/


package utils

import (
    "time"
    "strings"
    "fmt"
)

var conversion = map[string]string{
    /*stdLongMonth      */ "B": "January",
    /*stdMonth          */ "b": "Jan",
    // stdNumMonth       */ "m": "1",
    /*stdZeroMonth      */ "m": "01",
    /*stdLongWeekDay    */ "A": "Monday",
    /*stdWeekDay        */ "a": "Mon",
    // stdDay            */ "d": "2",
    // stdUnderDay       */ "d": "_2",
    /*stdZeroDay        */ "d": "02",
    /*stdHour           */ "H": "15",
    // stdHour12         */ "I": "3",
    /*stdZeroHour12     */ "I": "03",
    // stdMinute         */ "M": "4",
    /*stdZeroMinute     */ "M": "04",
    // stdSecond         */ "S": "5",
    /*stdZeroSecond     */ "S": "05",
    /*stdLongYear       */ "Y": "2006",
    /*stdYear           */ "y": "06",
    /*stdPM             */ "p": "PM",
    // stdpm             */ "p": "pm",
    /*stdTZ             */ "Z": "MST",
    // stdISO8601TZ      */ "z": "Z0700",  // prints Z for UTC
    // stdISO8601ColonTZ */ "z": "Z07:00", // prints Z for UTC
    /*stdNumTZ          */ "z": "-0700", // always numeric
    // stdNumShortTZ     */ "b": "-07",    // always numeric
    // stdNumColonTZ     */ "b": "-07:00", // always numeric
}

/**
 * @desc    将时间格式化成需要的格式
 * @name    TimeFormat
 * @param   string      format
 * @param   time.Time   t
 * @return  string
 */
func TimeFormat(format string, t time.Time) string {
    formatChunks := strings.Split(format, "%")
    var layout []string
    for _, chunk := range formatChunks {
        if len(chunk) == 0 {
            continue
        }
        if layoutCmd, ok := conversion[chunk[0:1]]; ok {
            layout = append(layout, layoutCmd)
            if len(chunk) > 1 {
                layout = append(layout, chunk[1:])
            }
        } else {
            layout = append(layout, "%", chunk)
        }
    }
    return t.Format(strings.Join(layout, ""))
}

/**
 * @desc 格式化无序
 */
func LogNsdFormat(logmap map[string]string, separator string) string {
    var formats []string
    for key, value := range logmap {
        switch length := len(separator); length {
        case 2:
            tmp := fmt.Sprintf("%s%v%s%v", key, string(separator[0]), value, string(separator[1]))
            formats = append(formats, tmp)
        default:
            tmp := []string{key, value}
            formats = append(formats, strings.Join(tmp, separator))
        }
    }
    return strings.Join(formats, " ")
}

/**
 * @desc 格式化有序
 */
func LogSdFormat(logmap [][2]string, separator string) string {
    var formats []string
    for _, v := range logmap {
        switch length := len(separator); length {
        case 2:
            tmp := fmt.Sprintf("%s%v%s%v", v[0], string(separator[0]), v[1], string(separator[1]))
            formats = append(formats, tmp)
        default:
            tmp := fmt.Sprintf("%s%v%s", v[0], string(separator), v[1])
            formats = append(formats, tmp)
        }
    }
    return strings.Join(formats, " ")
}
