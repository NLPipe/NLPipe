#!/bin/bash

args=("$@")

file_name="${args[0]}"

min=`awk 'BEGIN{a=999999999999}{if ($1<0+a) a=$1} END{print a}' $file_name`
max=`awk 'BEGIN{a=           0}{if ($1>0+a) a=$1} END{print a}' $file_name`
avg=`awk '{ sum += $1; n++ } END { if (n > 0) print sum / n; }' $file_name`

echo "${min}, ${max}, ${avg}"
