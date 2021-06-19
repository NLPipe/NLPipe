import json
import os
import boto3
from decimal import Decimal

def lambda_handler(event, context):

    sentence = event['responsePayload']['sentence']
    file_name = event['responsePayload']['file_name']
    sentiment_analysis_sentiment = event['responsePayload']['sentiment_analysis_sentiment']
    sentiment_analysis_sentiment_positive_probability = event['responsePayload']['sentiment_analysis_sentiment_positive_probability']
    sentiment_analysis_sentiment_negative_probability = event['responsePayload']['sentiment_analysis_sentiment_negative_probability']

    print("[DEBUG] The following data will be sent to DynamoDB:")
    print("    sentence: " + sentence)
    print("    file_name: " + file_name)
    print("    sentiment_analysis_sentiment: " + sentiment_analysis_sentiment)
    print(
        "    sentiment_analysis_sentiment_positive_probability: " +
        str(sentiment_analysis_sentiment_positive_probability)
    )
    print(
        "    sentiment_analysis_sentiment_negative_probability: " +
        str(sentiment_analysis_sentiment_negative_probability)
    )
    print("\n\n")

    print("[DEBUG] Opening a connection with DynamoDB...")
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('sentiment_analysis_results')
    print("[DEBUG] DynamoDB connection OK!")
    print("\n\n")

    print("[DEBUG] Writing to DynamoDB...")
    table.put_item(
        Item={
            'file_name' : file_name,
            'sentence' : sentence,
            'sentiment' : sentiment_analysis_sentiment,
            'sentiment_prob_neg' : str(sentiment_analysis_sentiment_negative_probability),
            'sentiment_prob_pos' : str(sentiment_analysis_sentiment_positive_probability)

        }
    )
    print("[DEBUG] DynamoDB write OK!")
    print("\n\n")

    return {
        'statusCode': 200,
        'body': json.dumps('DynamoDB write OK!')
    }
