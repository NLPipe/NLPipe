args=("$@")

# parsing args from command line
sentences_max_num=${args[0]}

if [[ -z "${sentences_max_num// }" ]]; then
  sentences_max_num=1
fi

mkdir -p ../deploy/Lambda/sentiment_analysis/inputs

echo "Removing old files in current ../deploy/Lambda/sentiment_analysis/inputs directory..."
rm -f ../deploy/Lambda/sentiment_analysis/inputs/*.json
echo "Old files removal DONE!"

echo "Creating input JSON files..."

i=1
for file_name in ../sentences/*.txt; do
  echo "JSON $i / $sentences_max_num populated w/ sentences in file $file_name"

  file_content=$(<$file_name)

  # echo $file_content

  json="{
    \"responsePayload\": {
      \"statusCode\": 200,
      \"body\": \"${file_content}\",
      \"file_name\": \"${i}\",
      \"recognized_text\": \"${file_content}\"
    }
  }"

  echo $json >> "../deploy/Lambda/sentiment_analysis/inputs/${i}.json"


  ((i=i+1))

  if (( $i > $sentences_max_num )); then
    break
  fi
done

echo "JSON files creation DONE!"
