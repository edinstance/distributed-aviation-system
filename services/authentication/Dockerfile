# ---------- 1️⃣ Build stage ----------
FROM python:3.12-slim AS builder

ENV PYTHONDONTWRITEBYTECODE=1 PYTHONUNBUFFERED=1
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential libpq-dev \
 && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .

RUN pip install --upgrade pip \
 && pip install --no-cache-dir -r requirements.txt

COPY . .

FROM python:3.12-slim AS web-base

ENV PYTHONDONTWRITEBYTECODE=1 PYTHONUNBUFFERED=1 DJANGO_SETTINGS_MODULE=authentication.settings

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends libpq5 \
 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local /usr/local
COPY --from=builder /app /app

ENV PORT=8080

EXPOSE ${PORT}

FROM web-base AS web
CMD ["gunicorn", "--bind", "0.0.0.0:8000", "authentication.wsgi:application"]

FROM web-base AS migrate
CMD ["python", "manage.py", "migrate_schemas", "--shared", "--noinput"]