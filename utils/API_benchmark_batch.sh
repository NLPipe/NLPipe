#!/bin/bash

args=("$@")

# parsing args from command line
batch_size=${args[0]}
bool_verbose=${args[1]}

if [[ -z "${bool_verbose// }" ]]; then
  bool_verbose=0
fi

# computing some utility variables
process_pid=$$
output_file_name="\
  ../benchmark_results/API/batch_size_${batch_size}.txt\
"
touch $output_file_name
tmp_file_name="../benchmark_results/API/${process_pid}.txt"

# placing $batch_size """parallel""" speech2text calls
if [ "$bool_verbose" -eq "1" ]; then
  echo "Starting $batch_size batches..."
fi
for (( i=1; i<=$batch_size; i++ ))
do

  newman run ../deploy/API/inputs/apitest.json | grep "average response" >> ${tmp_file_name} &

done
if [ "$bool_verbose" -eq "1" ]; then
  echo "$batch_size batches started!"
fi

# waiting for all the """parallel""" speech2text to complete
if [ "$bool_verbose" -eq "1" ]; then
  echo "Waiting for $batch_size batches to finish..."
fi
wait
if [ "$bool_verbose" -eq "1" ]; then
  echo "$batch_size batches DONE!"
fi

cat $tmp_file_name

# min=`awk 'BEGIN{a=999999999999}{if ($1<0+a) a=$1} END{print a}' $tmp_file_name`
# max=`awk 'BEGIN{a=           0}{if ($1>0+a) a=$1} END{print a}' $tmp_file_name`
# avg=`awk '{ sum += $1; n++ } END { if (n > 0) print sum / n; }' $tmp_file_name`

# echo "${min}, ${max}, ${avg}"
# echo "${min}, ${max}, ${avg}" >> $output_file_name

rm -f $tmp_file_name
