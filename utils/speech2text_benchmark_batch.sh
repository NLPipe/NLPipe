#!/bin/bash

args=("$@")

# parsing args from command line
batch_size=${args[0]}
bool_verbose=${args[1]}
full_file_name="${args[2]}"
num_of_batches=${args[3]} # can be left empty!

if [[ -z "${num_of_batches// }" ]]; then
  num_of_batches=1
fi

# computing some utility variables
file_extension="${full_file_name##*.}"
file_name="${full_file_name%.*}"
dir_name=$file_name
process_pid=$$
output_file_name="\
  ../benchmark_results/speech2text/${dir_name}/${file_name}_${batch_size}_${num_of_batches}.txt\
"
tmp_file_name="../benchmark_results/speech2text/${dir_name}/${process_pid}.txt"

# placing $batch_size """parallel""" speech2text calls
if [ "$bool_verbose" -eq "1" ]; then
  echo "Starting $batch_size batches..."
fi
for (( i=1; i<=$batch_size; i++ ))
do

  iteration_file_name="${dir_name}/${file_name}_${i}.${file_extension}"

  python3 ../src/speech2text.py ../audio/"${iteration_file_name}" >> ${tmp_file_name} &

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

min=`awk 'BEGIN{a=999999999999}{if ($1<0+a) a=$1} END{print a}' $tmp_file_name`
max=`awk 'BEGIN{a=           0}{if ($1>0+a) a=$1} END{print a}' $tmp_file_name`
avg=`awk '{ sum += $1; n++ } END { if (n > 0) print sum / n; }' $tmp_file_name`

#echo "${min}, ${max}, ${avg}"
echo "${min}, ${max}, ${avg}" >> $output_file_name

rm -f $tmp_file_name
