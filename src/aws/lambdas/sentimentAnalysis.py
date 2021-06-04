import json

import nltk
import nltk.classify.util
from nltk.classify import NaiveBayesClassifier
from nltk.corpus import movie_reviews

nltk.data.path.append("/tmp")
nltk.download("movie_reviews", download_dir = "/tmp")

import string

import pickle

import sys

import boto3

def lambda_handler(event, context):

    # classify a sample input sentence to positive / negative class and give the probability of
    # belonging to each class, based on trained model

    def classify_text(text):
        tokenized_text = tokenizer(text)
        featured_text = word_feats(tokenized_text)
        result = classifier_from_disk.classify(featured_text)
        result_prob = classifier_from_disk.prob_classify(featured_text)
        return result, result_prob.prob('pos'), result_prob.prob('neg')


    # converting words inside input sentence ( set of words ) to a dictionary of [word , true] to set
    # learning features

    def word_feats(words):
        return dict([(word, True) for word in words])


    # tokenize ( separate ) words of input text, leave out punctuations ( . , ? ! : " ' etc )
    # and change to lower case format

    def tokenizer(text):
        stops = list(string.punctuation)
        tokens = []
        for word in text:
            word.lower()
            if word not in stops:
                tokens.append(word)
        return tokens


    s3 = boto3.resource('s3')
    classifier_from_disk = pickle.loads(
        s3.Bucket("nlpipe-test-stt")
          .Object("NaiveBayesClassifierTrainedModel.pickle").get()['Body'].read()
    )

    #result, result_prob_pos, result_prob_neg = classify_text(input_text)
    result, result_prob_pos, result_prob_neg = classify_text(
        "love emotion joy passion sun flower choccolate pizza pasta"
    )
    print("[DEBUG] Your sentence belong to class:", result)
    print("[DEBUG] probability of belong to class positive is :", result_prob_pos)
    print("[DEBUG] probability of belong to class negative is :", result_prob_neg)


    return {
        'statusCode': 200,
        'body': json.dumps('[DEBUG] Success!')
    }
