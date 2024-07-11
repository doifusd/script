#!/bin/bash
logPath="/Users/sky/Documents/log/result/"
for line in $(cat /Users/sky/Desktop/invoiceNo.txt)
do
 cat ${logPath}$1|grep $line >> ${logPath}$2
done
