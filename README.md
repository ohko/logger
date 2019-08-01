[![github.com/ohko/logger](https://goreportcard.com/badge/github.com/ohko/logger)](https://goreportcard.com/report/github.com/ohko/logger)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/ab57f8d1f67b47699af16eafc089f8bf)](https://www.codacy.com/app/ohko/logger?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ohko/logger&amp;utm_campaign=Badge_Grade)

# 日志打印管理
通过环境变量`LOG_LEVEL`可控制日志的输出等级。

```golang
// 默认仅显示在os.Stdout
ll := NewLogger(nil)

// 指定输出设备为os.Stdout
ll := NewLogger(os.Stdout)

// 使用内置的按日切割输出到文件
ll := NewLogger(NewDefaultWriter(nil))

// 内置的基础上同时显示在os.Stdout
ll := NewLogger(NewDefaultWriter(os.Stdout))

ll.SetLevel(LoggerLevel0Debug)
ll.Log0Debug(fmt.Sprintf("0:%v", "Debug"))
ll.Log1Warn("1:Warning")
ll.Log2Error("2:Error")
ll.Log3Fatal("3:Fatal") // 附加 os.Exit(1)
ll.Log4Trace("4:Trace")
```