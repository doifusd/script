#!/bin/bash
# cat $1 |grep -Eo 'invoiceNo:\":\"[0-9]+' > "/Users/sky/Desktop/${2}.log"
cat $1 |grep -Eo 'invoiceNo:\":\"[0-9]+'
