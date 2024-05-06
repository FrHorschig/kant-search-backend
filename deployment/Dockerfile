FROM golang:1.22
WORKDIR /app
COPY src/kant-search-backend .

RUN apt-get update
RUN apt-get install -y python3 python3-venv

RUN python3 -m venv .venv
RUN source .venv/bin/activate
RUN pip install -U pip setuptools wheel
RUN pip install -U spacy
RUN python -m spacy download en_core_web_sm
RUN deactivate

COPY src_py/split_text.py .
ENV KSGO_PYTHON_BIN_PATH=".venv/bin/python3"
ENV KSGO_PYTHON_SCRIPT_PATH="split_text.py"

CMD ["./kant-search-backend"]