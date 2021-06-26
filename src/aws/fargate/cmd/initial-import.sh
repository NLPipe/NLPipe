aws dynamodb create-table --cli-input-json file://$(dirname "$0")/table.json --endpoint-url http://127.0.0.1:8000
aws dynamodb batch-write-item --request-items file://$(dirname "$0")/items.json --endpoint-url http://127.0.0.1:8000
aws dynamodb scan --table-name=nlpipe-results --endpoint-url http://localhost:8000
