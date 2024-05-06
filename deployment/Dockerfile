FROM golang:1.22
WORKDIR /app
COPY src/kant-search-backend .

RUN apt-get update && apt-get install -y python3 python3-pip
RUN pip3 install spacy
COPY src_py/split_text.py .
ENV KSGO_PYTHON_BIN_PATH="python3"
ENV KSGO_PYTHON_SCRIPT_PATH="split_text.py"

CMD ["./kant-search-backend"]