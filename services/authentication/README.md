# Authentication Service

Multi-tenant Django authentication service for the distributed aviation system.

## Features

- Multi-tenant architecture using django-tenants
- JWT authentication with Django REST Framework
- PostgreSQL database backend
- CORS support for frontend integration
- Organization-based tenant isolation

## Setup

### Prerequisites
- Python 3.12+
- PostgreSQL
- Virtual environment

### Installation

1. **Create and activate virtual environment:**
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

2. **Install dependencies:**
```bash
pip install -r requirements.txt
```

3. **Configure environment:**
```bash
cp .env.example .env
# Edit .env with your settings
```

4. **Generate RSA keys for JWT:**
```bash
python keys/generate_keys.py
```

5. **Database setup:**
```bash
# Create PostgreSQL database
createdb -U postgres authentication

# Run migrations
python manage.py migrate_schemas --shared
```

## Usage

### Running the Server
```bash
python manage.py runserver
```

### API Access
Most APIs will require the ```X-Org-id``` header.


## Configuration

### Environment Variables
- `SECRET_KEY`: Django secret key
- `DEBUG`: Debug mode (True/False)
- `ALLOWED_HOSTS`: Comma-separated list of allowed hosts
- `DATABASE_URL`: PostgreSQL connection string
- `KEYS_DIR`: Directory containing RSA keys (`./keys` for local, `/keys` for Docker)

### JWT Keys Setup

The authentication service uses RSA keys for JWT token signing and verification.

**Local Development:**
- Keys are stored in `./keys/` directory
- Run `python keys/generate_keys.py` to generate keys
- Set `KEYS_DIR=./keys` in your `.env` file

**Docker Environment:**
- Keys are generated automatically by the `authentication-keygen` service
- Keys are shared between containers via the `jwt_keys` volume
- Set `KEYS_DIR=/keys` in Docker environment

**Key Files:**
- `private.pem`: RSA private key for signing JWTs
- `public.pem`: RSA public key for verifying JWTs
- `keymap.json`: Key metadata for rotation support

### Multi-Tenant Setup
- **Shared Apps**: Organizations, core Django apps
- **Tenant Apps**: Users, authentication API, admin
- **Database Router**: Automatically routes queries to correct schema

## API Endpoints

### Authentication
- `POST /api/auth/login/` - User login
- `POST /api/auth/logout/` - User logout
- `POST /api/auth/refresh/` - Refresh JWT token

### Users
- `POST /api/users/create` - Create user

## Database Schema

### Public Schema (Shared)
- Organizations
- Domains
- Django core tables

### Tenant Schema (Per Organization)
- Users
- Sessions
- Admin logs
- API tokens

## Development

### Running Tests
```bash
python manage.py test
```

### Code Quality
```bash
# Check for issues
python manage.py check

# Validate migrations
python manage.py showmigrations
```