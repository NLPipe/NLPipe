import json
import urllib.parse
import boto3

import sys
import speech_recognition as sr

s3 = boto3.client('s3')

def lambda_handler(event, context):

    #bucket = event['Records'][0]['s3']['bucket']['name']
    bucket = event['name']
    key = urllib.parse.unquote_plus(
        #event['Records'][0]['s3']['object']['key'], encoding='utf-8'
        event['key'], encoding='utf-8'
    )

    try:
        ### --- REMEMBER TO CHANGE BUCKET NAME and FILE --- ###
        response = s3.get_object(Bucket='nlpipe-test-stt', Key='a.wav')
        ### --- REMEMBER TO CHANGE BUCKET NAME and FILE --- ###

        #print("[DEBUG] CONTENT TYPE: " + str(response.keys()))
        #print("[DEBUG] BODY: " + str(response['Body']))

    except Exception as e:
        print(e)
        print(
            'Error getting object {} from bucket {}' +
            'Make sure they exist and your bucket is in the same region as this function.'
            .format(key, bucket)
        )
        raise e

    recognizer = sr.Recognizer()


    with sr.AudioFile(response['Body']) as source:
        print('[DEBUG] Opened file')
        recorded_audio = recognizer.listen(source)

    try:
        text = recognizer.recognize_google(recorded_audio, language="en-US")
        print('[DEBUG] recognized text: ' + text)

    except Exception as ex:
        print(ex)

    return {
        'statusCode': 200,
        'body': json.dumps('Success!')
    }
