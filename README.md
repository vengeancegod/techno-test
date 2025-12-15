CLI Task manager
Консольное приложение для создания и управленя задачами, разработанное на go с использование postgreSQL, docker

Архитектура в проекте представлена следующим образом:
techno-test/
├── cmd/
│   ├── app/              # Основное приложение (API)
│   │   └── main.go
│   └── worker/           # Worker приложение
│       └── main.go
├── internal/
│   ├── app/             # Логика приложения
│   │   ├── app.go       # Основное приложение
│   │   ├── worker.go    # Worker приложение (таймер для очистки задач)
│   │   └── service_provider.go
│   ├── config/          # Конфигурация
│   ├── model/          # Сущности БД
│   ├── cli/            # CLI
│   ├── repository/      # Репозитории (DB доступ)
│   └── service/         # Бизнес-логика
├── migrations/          # Миграции БД
├── .env.example        # Пример переменных окружения
├── .gitignore
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile            # Утилиты для сборки и запуска
└── README.md