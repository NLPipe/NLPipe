# invokes speech2text AWS lambda function.
# keep --invocation-type as it is
# --payload is a JSON file storing the content of the JSON test file set in AWS
# speech2text_output is a JSON file storing the output of the CLI command (usually, just a code 202 response)

args=("$@")
batch_size=${args[0]}
# full_file_name="${args[1]}"

# file_name="${full_file_name%.*}"

echo "Placing $batch_size parallel calls..."

for (( i=1; i<=$batch_size; i++ ))
do

  aws lambda invoke \
  --function-name sentiment_analysis \
  --invocation-type Event \
  --payload "fileb://../deploy/Lambda/sentiment_analysis/inputs/${i}.json" \
  ../deploy/Lambda/sentiment_analysis/outputs/sentiment_analysis_output.json >> \
  ../deploy/Lambda/sentiment_analysis/outputs/sentiment_analysis_output.txt \
  &

done

echo ""
current_time=$(date +"%H:%M")
echo "$batch_size parallel calls placed! --> $current_time"

echo "Waiting for $batch_size parallel calls to be completed..."
wait
echo "$batch_size parallel calls DONE!"
