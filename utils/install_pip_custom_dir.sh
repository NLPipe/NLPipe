mkdir lambda_layers && cd lambda_layers
mkdir python && cd python
pip3 install requests -t ./
cd .. && zip -r python_modules.zip .
