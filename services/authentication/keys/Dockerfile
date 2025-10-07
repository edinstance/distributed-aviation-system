FROM python:3.12-alpine

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY generate_keys.py .

CMD ["python", "generate_keys.py"]