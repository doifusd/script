#!/bin/bash
# cat /Users/sky/Desktop/zaiciqingqiu.txt | while read line
# do
# grep $line /Users/sky/Documents/script/log/result/* >> ./3.txt
#
# done

keyword=$(cat /Users/sky/Desktop/zaiciqingqiu.txt)

# 在目录B中查找包含关键词的文件
# find /Users/sky/Documents/script/log/result/ -name "*\.log" -exec grep -H "$keyword" {} \;
find /Users/sky/Documents/script/log/result/ -name "*\.log" -exec grep -H "$keyword" {} \;
