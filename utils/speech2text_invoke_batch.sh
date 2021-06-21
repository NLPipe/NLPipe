# invokes speech2text AWS lambda function.
# keep --invocation-type as it is
# --payload is a JSON file storing the content of the JSON test file set in AWS
# speech2text_output is a JSON file storing the output of the CLI command (usually, just a code 202 response)

args=("$@")
batch_size=${args[0]}
file_name="${args[1]}"

echo "Placing $batch_size parallel calls..."

for (( i=1; i<=$batch_size; i++ ))
do

  aws lambda invoke \
  --function-name speech2text \
  --invocation-type Event \
  --payload "fileb://../deploy/Lambda/speech2text/inputs/${file_name}/${file_name}_${i}.json" \
  ../deploy/Lambda/speech2text/outputs/speech2text_output.json > \
  ../deploy/Lambda/speech2text/outputs/speech2text_output.txt &

done

echo "$batch_size parallel calls placed!"

echo "Waiting for $batch_size parallel calls to be completed..."
wait
echo "$batch_size parallel calls DONE!"
