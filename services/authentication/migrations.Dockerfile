FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    DJANGO_SETTINGS_MODULE=authentication.settings

WORKDIR /app

# Install minimal OS libs for PostgreSQL
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential libpq-dev \
 && rm -rf /var/lib/apt/lists/*
    # Copy and install only migration requirements
COPY migration.requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Copy only necessary files for migrations
COPY . .

# Default command: run tenant schemas migrations
CMD ["python", "manage.py", "migrate_schemas", "--tenant", "--shared", "--noinput"]