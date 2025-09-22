# 🌿 Vitalis Life - Серверная часть


Серверное приложение для интернет-магазина здорового питания и экопродуктов Vitalis Life.

## 🚀 Возможности
1. RESTful API для фронтенд-приложения

2. Обработка заказов и платежей через ЮКассу

3. Интеграция с SMTP для email-уведомлений

4. Логирование и CORS поддержка

📋 Предварительные требования
Go 1.19 или выше

PostgreSQL 12+

Доступ к SMTP-серверу

Аккаунт в ЮКассе (для обработки платежей)

⚙️ Настройка конфигурации
1. Создание файла конфигурации
Создайте файл config.yaml в папке config/ со следующим содержимым:

yaml
server:
  version: "Vitalis-api-v1.0.0"
  port: ":8080"

logger:
  level: "debug" # debug, info

cors:
  allow_origins:
    - "*"
  allow_methods:
    - "GET"
    - "POST"
    - "PATCH"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allow_headers:
    - "Origin"
    - "Content-Type"
    - "Accept"
  expose_headers: []
  allow_credentials: false
  max_age: 43200 # 12 часов в секундах
2. Настройка переменных окружения
Создайте файл .env в папке cmd/ или экспортируйте переменные в среде выполнения:

env
# Настройки базы данных PostgreSQL
DB_USER=your_database_user
DB_PASSWORD=your_secure_password
DB_NAME=vitalis_db
DB_HOST=localhost
DB_PORT=5432

# Настройки SMTP сервера для email-уведомлений
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# Настройки платежной системы ЮКасса
YOOKASSA_SHOP_ID=your_shop_id
YOOKASSA_SECRET_KEY=your_secret_key

# Доменное имя фронтенд-приложения
FRONTEND_URL=https://vitalis-life.ru

# Email менеджера для уведомлений о заказах
MANAGER_EMAIL=manager@vitalis-life.ru
🔧 Описание параметров окружения
Настройки базы данных
DB_USER - имя пользователя PostgreSQL

DB_PASSWORD - пароль пользователя базы данных

DB_NAME - название базы данных

DB_HOST - хост базы данных

DB_PORT - порт для подключения к БД

Настройки email-рассылок
SMTP_HOST - адрес SMTP-сервера для отправки писем

SMTP_PORT - порт SMTP-сервера

SMTP_USER - логин для авторизации на SMTP-сервере

SMTP_PASSWORD - пароль для авторизации на SMTP-сервере

Настройки платежной системы
YOOKASSA_SHOP_ID - идентификатор магазина в ЮКассе

YOOKASSA_SECRET_KEY - секретный ключ для доступа к API ЮКассы

Общие настройки
FRONTEND_URL - базовый URL фронтенд-приложения

MANAGER_EMAIL - email адрес менеджера для получения уведомлений о заказах

🚀 Запуск приложения
Клонируйте репозиторий:

bash
git clone <repository-url>
cd vitalis-server
Установите зависимости:

bash
go mod download
Настройте конфигурацию (как описано выше)

Запустите приложение:

bash
go run cmd/main.go
Приложение будет доступно по адресу: http://localhost:8080

📁 Структура проекта
text
vitalis-server/
├── cmd/                 # Точка входа приложения
├── config/              # Конфигурационные файлы
├── internal/            # Внутренние пакеты приложения
│   ├── controller/      # HTTP контроллеры
│   ├── service/         # Бизнес-логика
│   ├── repository/      # Работа с данными
│   └── model/           # Структуры данных
├── pkg/                 # Вспомогательные пакеты
└── docs/                # Документация
📄 Лицензия
Этот проект является собственностью Vitalis Life. Все права защищены.
