#!/bin/bash

args=("$@")

file_name="${args[0]}"

repeated_min=`awk '{ sum += $1; n++ } END { if (n > 0) print sum / n; }' $file_name`
repeated_max=`awk '{ sum += $2; n++ } END { if (n > 0) print sum / n; }' $file_name`
repeated_avg=`awk '{ sum += $3; n++ } END { if (n > 0) print sum / n; }' $file_name`

echo "${repeated_min}, ${repeated_max}, ${repeated_avg}"
