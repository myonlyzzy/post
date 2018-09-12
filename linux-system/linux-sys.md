# linux system

记录一些linux下系统操作的相关命令

## linux下查看进程的io读写情况

1. iotop -oP  o表示只显示当前有io读写的进程 P表示显示进程信息
2. pidstat -d 1 展示I/O统计，每秒更新一次
3. wc -l 用来统计行数 例如 netstat -an |grep 3306  |wc -l  统计连接数
   pstree 123 |wc -l 统计线程数
4. 输出一列 awk '{print $1}'

## 其他常用命令

1. Ulimit -a 查看系统资源的限制,例如可以打开的最大文件描述符数量

