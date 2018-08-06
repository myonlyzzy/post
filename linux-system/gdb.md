# linux 下gdb的使用


## 使用gdb 调试golang程序

### 1. attach 进golang进程
1. 编译golang代码
    ```
        go build -gcflags "-N -l" -o test main.go
    ```
    这里使用gcflags 是为了让golang 编译器关闭对内联函数的优化,方便gdb识别代码.
2.  启动编译好的golang 执行文件

3.  gdb attach pid ,pid 为goalng程序的进程号.
    info threads 显示所有的线程
    thread id , 切换到id号表示的线程
    bt ,当前线程的调用栈
    thread apply all bt ,显示所有的线程调用栈




