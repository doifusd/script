#!/bin/bash
cat /Users/sky/Desktop/2023-08-10/tmp.txt | while read line
do
    ./jsonClean -str "$line" >> $1
done
