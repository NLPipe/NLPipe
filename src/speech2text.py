import sys
import speech_recognition as sr
import datetime

recognizer = sr.Recognizer()

start_time = datetime.datetime.now()

with sr.AudioFile(sys.argv[1]) as source:
    recorded_audio = recognizer.listen(source)

try:
    text = recognizer.recognize_google(
            recorded_audio,
            language="en-US"
        )
    # print(text)

except Exception as ex:
    # print(ex)
    x=1

end_time = datetime.datetime.now()

time_diff = (end_time - start_time)
execution_time = time_diff.total_seconds() * 1000

# print(str(execution_time) + " ms")
print(str(execution_time))
