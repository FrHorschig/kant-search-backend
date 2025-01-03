FROM golang:1.22
WORKDIR /app
COPY src/kant-search-backend .

RUN apt-get update
RUN apt-get install -y python3 python3-venv

RUN python3 -m venv .venv
RUN . .venv/bin/activate \
  pip install -U setuptools wheel spacy \
  python -m spacy download en_core_web_sm
COPY src_py/split_text.py .
ENV KSGO_PYTHON_BIN_PATH=".venv/bin/python3"
ENV KSGO_PYTHON_SCRIPT_PATH="."

EXPOSE 3000
CMD ["./kant-search-backend"]