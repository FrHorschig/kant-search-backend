import sys
import json
from spacy.lang.de import German


def main():
    data = json.load(sys.stdin)
    nlp = German()
    nlp.add_pipe("sentencizer")

    result = {}
    for item in data:
        doc = nlp(item["Text"])
        sentences = [sent.text for sent in doc.sents]
        result[str(item["Id"])] = sentences

    print(json.dumps(result))


if __name__ == "__main__":
    main()
