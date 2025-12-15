CLI Task manager.
Консольное приложение для создания и управленя задачами, разработанное на go с использование postgreSQL, docker
```bash
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
```
Запуск приложения:
```bash
git clone https://github.com/vengeancegod/techno-test.git
cd techno-test

docker compose up -d --build

make build
```
Запуск утилиты:
```bash
./bin/taskmanager
```
Запуск воркера (тот, что удаляет выполненные задачи из БД):
```bash
./bin/taskcleaner
```

Создание миграций:
```bash
make migrate-up
```
Работа API:
Создание задачи
```bash
bin/taskmanager task create -t "Подготовить отчет" -d "Ежеквартальный отчет по продажам"
```
Получение задач:
```bash
bin/taskmanager task list
```
Фильтрация по статусу задач (done/not_done):
```bash
bin/taskmanager task list -s done
bin/taskmanager task list --status not_done
```
Получение задачи по ID:
```bash
bin/taskmanager task get 42
```
Обновление задачи:
```bash
bin/taskmanager task update 1 -t "Обновленный заголовок" -d "Обновленное описание"
```
Удаление задачи:
```bash
bin/taskmanager task delete
```

