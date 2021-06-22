docker build -t nlpipe .
docker tag nlpipe:latest $(echo $REPO_URL)
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $(echo $URL)
docker push $(echo $REPO_URL)
