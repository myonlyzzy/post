# linux system

记录一些linux下系统操作的相关命令

## linux下查看进程的io读写情况

1. iotop -oP  o表示只显示当前有io读写的进程 P表示显示进程信息

2. pidstat -d 1 展示I/O统计，每秒更新一次