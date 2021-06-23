args=("$@")
batch_size=${args[0]}
full_file_name="${args[1]}"

extension="${full_file_name##*.}"
file_name="${full_file_name%.*}"

mkdir -p ../deploy/Lambda/speech2text/inputs/

for (( i=1; i<=$batch_size; i++ ))
do
  json="{
    \"Records\": [
      {
        \"eventVersion\": \"2.0\",
        \"eventSource\": \"aws:s3\",
        \"awsRegion\": \"us-east-1\",
        \"eventTime\": \"1970-01-01T00:00:00.000Z\",
        \"eventName\": \"ObjectCreated:Put\",
        \"userIdentity\": {
          \"principalId\": \"EXAMPLE\"
        },
        \"requestParameters\": {
          \"sourceIPAddress\": \"127.0.0.1\"
        },
        \"responseElements\": {
          \"x-amz-request-id\": \"EXAMPLE123456789\",
          \"x-amz-id-2\": \"EXAMPLE123/5678abcdefghijklambdaisawesome/mnopqrstuvwxyzABCDEFGH\"
        },
        \"s3\": {
          \"s3SchemaVersion\": \"1.0\",
          \"configurationId\": \"testConfigRule\",
          \"bucket\": {
            \"name\": \"nlpipe-test-stt\",
            \"ownerIdentity\": {
              \"principalId\": \"EXAMPLE\"
            },
            \"arn\": \"arn:aws:s3:::nlpipe-test-stt\"
          },
          \"object\": {
            \"key\": \"${file_name}/${file_name}_${i}.${extension}\",
            \"size\": 1024,
            \"eTag\": \"0123456789abcdef0123456789abcdef\",
            \"sequencer\": \"0A1B2C3D4E5F678901\"
          }
        }
      }
    ]
  }"

  echo $json > ../deploy/Lambda/speech2text/inputs/"${file_name}"/"${file_name}_${i}.json"

done
