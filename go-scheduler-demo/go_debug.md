#  GODEBUG 的使用

> 这篇文章不是介绍怎么debug golang程序,而是使用GODEBUG参数追踪程序里面的调度变化.

## 如何使用GODEBUG

```
var debug struct {
	allocfreetrace   int32
	cgocheck         int32
	efence           int32
	gccheckmark      int32
	gcpacertrace     int32
	gcshrinkstackoff int32
	gcrescanstacks   int32
	gcstoptheworld   int32
	gctrace          int32
	invalidptr       int32
	sbrk             int32
	scavenge         int32
	scheddetail      int32
	schedtrace       int32
}

var dbgvars = []dbgVar{
	{"allocfreetrace", &debug.allocfreetrace},
	{"cgocheck", &debug.cgocheck},
	{"efence", &debug.efence},
	{"gccheckmark", &debug.gccheckmark},
	{"gcpacertrace", &debug.gcpacertrace},
	{"gcshrinkstackoff", &debug.gcshrinkstackoff},
	{"gcrescanstacks", &debug.gcrescanstacks},
	{"gcstoptheworld", &debug.gcstoptheworld},
	{"gctrace", &debug.gctrace},
	{"invalidptr", &debug.invalidptr},
	{"sbrk", &debug.sbrk},
	{"scavenge", &debug.scavenge},
	{"scheddetail", &debug.scheddetail},
	{"schedtrace", &debug.schedtrace},
}
```
在runtime包的runtime1.go 文件中定义了GODEBUG的参数信息,如果要使用某个参数就写入环境变量中。例如,GODEBUG=schedtrace=2000,scheddetail=1.


## schedtrace scheddetail 详解

schedtrace schedetail 可以用来追踪调度相关的信息

实际的调度代码在runtime的 proc.go 文件中,sysmon函数调用schedtrace来输出信息

```
SCHED 918695ms: gomaxprocs=4 idleprocs=4 threads=8 spinningthreads=0 idlethreads=4 runqueue=0 gcwaiting=0 nmidlelocked=0 stopwait=0 sysmonwait=0
  P0: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0
  P1: status=0 schedtick=2 syscalltick=0 m=-1 runqsize=0 gfreecnt=0
  P2: status=0 schedtick=2 syscalltick=30 m=-1 runqsize=0 gfreecnt=0
  P3: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0
  M7: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 helpgc=0 spinning=false blocked=true lockedg=-1
  M6: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 helpgc=0 spinning=false blocked=true lockedg=-1
  M5: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 helpgc=0 spinning=false blocked=true lockedg=-1
  M4: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 helpgc=0 spinning=false blocked=true lockedg=-1
  M3: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 helpgc=0 spinning=false blocked=true lockedg=-1
  M2: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 helpgc=0 spinning=false blocked=false lockedg=-1
  M1: p=-1 curg=17 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 helpgc=0 spinning=false blocked=false lockedg=17
  M0: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 helpgc=0 spinning=false blocked=false lockedg=-1
  G1: status=4(IO wait) m=-1 lockedm=-1
  G17: status=6() m=1 lockedm=1
  G2: status=4(force gc (idle)) m=-1 lockedm=-1
  G3: status=4(GC sweep wait) m=-1 lockedm=-1
  G18: status=4(finalizer wait) m=-1 lockedm=-1
```
上面这段是打印出关于调度追踪的信息.
* gomaxprocs  //最大的procs信息,默认为	runtime.NumCPU()，cpu核数
* idleprocs   //空闲的processor
* threads     //创建的线程数,即machine数量
* spinningthread  //自旋状态的m数量
* idlethread  //空闲状态的线程



