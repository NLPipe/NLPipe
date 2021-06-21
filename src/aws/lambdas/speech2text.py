import json
import urllib.parse
import boto3

import sys

import speech_recognition as sr

s3 = boto3.client('s3')


def lambda_handler(event, conrecognized_text):
    #print("Received event: " + json.dumps(event, indent=2))

    # Get the object from the event and show its content type
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = urllib.parse.unquote_plus(event['Records'][0]['s3']['object']['key'], encoding='utf-8')

    try:
        response = s3.get_object(Bucket=bucket, Key=key)
        print("[DEBUG] Content type: " + response['ContentType'] + "\n\n")
        print("[DEBUG] Bucket      : " + bucket + "\n\n")
        print("[DEBUG] file_name    : " + key + "\n\n")
        #return response['ContentType']
    except Exception as e:
        print(e)
        print('\n\nError getting object {} from bucket {}. Make sure they exist and your bucket is in the same region as this function.\n\n'.format(key, bucket))
        raise e

    print('[DEBUG] Getting a sound recognizer instance\n\n')

    recognizer = sr.Recognizer()

    print('[DEBUG] Got     a sound recognizer instance\n\n')

    with sr.AudioFile(response['Body']) as source:
        print('[DEBUG] Opened file ' + key + " from bucket " + bucket + "\n\n")
        recorded_audio = recognizer.listen(source)

    try:
        recognized_text = recognizer.recognize_google(recorded_audio, language="en-US")
        print('[DEBUG] Recognized text: ' + recognized_text + "\n\n")

    except Exception as ex:
        print(ex)


    #return response['ContentType']
    #return recognized_text
    return {
            'statusCode': 200,
            'body': json.dumps(recognized_text),
            'file_name': key,
            'recognized_text': recognized_text
        }
