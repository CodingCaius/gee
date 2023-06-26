
/*
Go 语言中，比较常见的错误处理方法是返回 error，由调用者决定后续如何处理。但是如果是无法恢复的错误，可以手动触发 panic，当然如果在程序运行过程中出现了类似于数组越界的错误，panic 也会被触发。panic 会中止当前执行的程序，退出
panic 会导致程序被中止，但是在退出前，会先处理完当前协程上已经defer 的任务，执行完成后再退出

可以 defer 多个任务，在同一个函数中 defer 多个任务，会逆序执行。即先执行最后 defer 的任务。
在这里，defer 的任务执行完成之后，panic 还会继续被抛出，导致程序非正常结束

//在Go语言中，recover函数用于捕获并处理发生在函数调用过程中的panic异常，避免程序因为未处理的异常而崩溃

Go 语言还提供了 recover 函数，可以避免因为 panic 发生而导致整个程序终止，recover 函数只在 defer 中生效

错误处理也可以作为一个中间件，增强 gee 框架的能力
*/

package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

//用于打印堆栈跟踪信息以进行调试
//函数trace的输入参数是一个字符串message，它表示需要在打印的堆栈跟踪信息前添加的附加消息
/*
1, 创建一个长度为32的uintptr类型的数组pcs，用于存储程序计数器（program counter）的地址。
2, 使用runtime.Callers函数获取当前调用栈的信息，并将调用栈的信息填充到pcs数组中。runtime.Callers的第一个参数是要跳过的调用层数，这里传递了3表示跳过前三层调用，以避免打印trace函数和recover函数自身的调用信息。
3, 创建一个strings.Builder类型的变量str，用于构建最终的堆栈跟踪信息。
4, 将传入的message添加到str中，并在其后添加字符串"\nTraceback:"，用于标识跟踪信息的开始。
5, 遍历pcs数组中的地址，对每个地址执行以下操作：
使用runtime.FuncForPC函数获取与给定程序计数器地址关联的函数。
使用fn.FileLine函数获取函数对应的文件名和行号。
使用fmt.Sprintf将文件名和行号格式化为字符串，并将其添加到str中。
6, 返回最终构建的堆栈跟踪信息，即str的字符串表示。
*/
func trace(message string) string {
    var pcs [32]uintptr
    n := runtime.Callers(3, pcs[:])

    var str strings.Builder
    str.WriteString(message + "\nTraceback:")
    for _, pc := range pcs[:n] {
        fn := runtime.FuncForPC(pc)
        file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
    }
    return str.String()
}


//用于处理恢复（recovery）的函数
func Recovery() HandlerFunc {
    return func(c *Context) {
        defer func() {
            //调用recover()函数捕获可能发生的panic异常
            if err := recover(); err != nil {
                message := fmt.Sprintf("%s", err)
                log.Printf("%s\n\n", trace(message))
                c.Fail(http.StatusInternalServerError, "Internal Server Error")
            }
        }()  //() 表示对匿名函数的立即调用
        c.Next()
    }
}



/*
在没有错误恢复功能的情况下，一旦程序发生 panic 异常，它会打印出相关的错误信息和堆栈跟踪，然后立即退出。这种行为是为了确保错误被及时发现并尽早修复，以避免继续执行潜在有问题的程序。
*/