# invokes speech2text AWS lambda function.
# keep --invocation-type as it is
# --payload is a JSON file storing the content of the JSON test file set in AWS
# speech2text_output is a JSON file storing the output of the CLI command (usually, just a code 202 response)

args=("$@")
file_name="${args[0]}"

aws lambda invoke \
--function-name speech2text \
--invocation-type Event \
--payload "fileb://../inputs/speech2text_input_${file_name}.json" \
../outputs/speech2text_output.json
