# PR Reviewer Service

Сервис для автоматического назначения ревьюеров на Pull Request'ы.

## Запуск

```bash
make build  
make up     
```

Перед запуском создайте `.env` файл в корне проекта:

```env
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=review_service
POSTGRES_HOST=db
POSTGRES_PORT=5432

SERVER_PORT=8080
SERVER_WRITE_TIMEOUT=15s
SERVER_READ_TIMEOUT=15s
SERVER_IDLE_TIMEOUT=60s

MIGRATIONS_PATH=db/migrations

DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_CONN_MAX_LIFETIME=5
```

## Проблемы

**Массовая деактивация пользователей**:  
При деактивации всех пользователей команды, открытые PR остаются с текущими ревьюерами. Новые назначения невозможны до появления активных пользователей.
