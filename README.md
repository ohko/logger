[![github.com/ohko/logger](https://goreportcard.com/badge/github.com/ohko/logger)](https://goreportcard.com/report/github.com/ohko/logger)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/ab57f8d1f67b47699af16eafc089f8bf)](https://www.codacy.com/app/ohko/logger?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ohko/logger&amp;utm_campaign=Badge_Grade)

# 日志打印管理
通过环境变量`LOG_LEVEL`可控制日志的输出等级。

```golang
ll := NewLogger()
l1.SetLevel(LoggerLevel0Debug)
ll.Log0Debug(fmt.Sprintf("0:%v", "Debug"))
ll.Log1Warn("1:Warning")
ll.Log2Error("2:Error")
ll.Log3Fatal("3:Fatal")
ll.Log4Trace("4:Trace")
```