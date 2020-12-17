FROM python:3.8-buster

RUN python -m pip install --upgrade pip
RUN pip install apache-beam[gcp,test]

COPY start.sh /
RUN chmod 777 /start.sh

ENTRYPOINT [ "/start.sh" ]