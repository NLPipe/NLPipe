import sys
import speech_recognition as sr

recognizer = sr.Recognizer()

with sr.AudioFile(sys.argv[1]) as source:
    recorded_audio = recognizer.listen(source)

try:
    text = recognizer.recognize_google(
            recorded_audio, 
            language="en-US"
        )
    print(text)

except Exception as ex:
    print(ex)

