#!/bin/bash

args=("$@")

# parsing args from command line
batch_size=${args[0]}
num_of_batches=${args[1]}
bool_verbose=${args[2]}
full_file_name="${args[3]}"

# computing some utility variables
file_extension="${full_file_name##*.}"
file_name="${full_file_name%.*}"
dir_name=$file_name
# process_pid=$(date +%s)
process_pid=$$
output_file_name="\
  ../benchmark_results/speech2text/${dir_name}/${file_name}_${batch_size}.txt\
"

echo "Batch size: $batch_size"

# placing $num_of_batches speech2text_benchmark_parallel calls
for (( i=1; i<=$num_of_batches; i++ ))
do

  bash speech2text_benchmark_batch.sh $batch_size 0 $full_file_name $num_of_batches

  if [ "$bool_verbose" -eq "1" ]; then
    echo "Batch $i/$num_of_batches DONE!"
  fi

done

echo ""
echo ""

repeated_min=`awk '{ sum += $1; n++ } END { if (n > 0) print sum / n; }' $output_file_name`
repeated_max=`awk '{ sum += $2; n++ } END { if (n > 0) print sum / n; }' $output_file_name`
repeated_avg=`awk '{ sum += $3; n++ } END { if (n > 0) print sum / n; }' $output_file_name`

echo "Updated min, max and avg: ${repeated_min}, ${repeated_max}, ${repeated_avg}"

# waiting for all the """parallel""" speech2text to complete
# wait
