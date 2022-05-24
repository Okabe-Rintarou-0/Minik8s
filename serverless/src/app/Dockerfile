FROM python:3.10-slim

COPY ./  /app

WORKDIR /app

RUN pip install -r ./requirements.txt

EXPOSE 8080

ENTRYPOINT ["python"]

CMD ["func.py"]
