import nltk.classify.util
from nltk.classify import NaiveBayesClassifier
from nltk.corpus import movie_reviews
import string

import nltk
nltk.download('movie_reviews')

import pickle

import sys
import speech_recognition as sr

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

f = open('NaiveBayesClassifierTrainedModel.pickle', 'rb')
classifier_from_disk = pickle.load(f)
f.close()

print('\nNaiveBayesClassifier loaded from disk!\n')

# input_text = input("Enter a sentence to classify: ")
# result, result_prob_pos, result_prob_neg = classify_text(input_text)
result, result_prob_pos, result_prob_neg = classify_text(text)
print("Your sentence belong to class:", result)
print("probability of belong to class positive is :", result_prob_pos)
print("probability of belong to class negative is :", result_prob_neg)
