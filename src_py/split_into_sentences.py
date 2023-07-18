import sys
from spacy.lang.de import German
import json


def main():
    text = sys.argv[1]
    nlp = German()
    nlp.add_pipe("sentencizer")
    doc = nlp(text)
    sentences = [sent.text for sent in doc.sents]
    print(json.dumps(sentences))


if __name__ == "__main__":
    main()
