mkdir sr
cd sr
python3 -m virtualenv venv
source venv/bin/activate
pip3 install speechrecognition
deactivate
cd venv/lib/python3.9/site-packages
zip -r ~/Desktop/sr/sr.zip speech_recognition
