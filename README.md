🌿 Vitalis Life - Серверная часть
<details> <summary>📋 Оглавление</summary>
🚀 Возможности

📋 Предварительные требования

⚙️ Настройка конфигурации

1. Создание файла конфигурации

2. Настройка переменных окружения

🔧 Описание параметров окружения

🗄️ Настройки базы данных

📧 Настройки email-рассылок

💳 Настройки платежной системы

🌐 Общие настройки

🚀 Запуск приложения

📁 Структура проекта

📄 Лицензия

</details>
Серверное приложение для интернет-магазина здорового питания и экопродуктов Vitalis Life.

🚀 Возможности
<kbd>REST API</kbd> <kbd>Аутентификация</kbd> <kbd>Платежи</kbd> <kbd>Email</kbd> <kbd>База данных</kbd>

<div align="center">
№	Функциональность	Описание
1	RESTful API	Полнофункциональное API для фронтенд-приложения
2	Управление пользователями	Регистрация, аутентификация и авторизация
3	Обработка платежей	Интеграция с ЮКассой для онлайн-платежей
4	Email-уведомления	Отправка писем через SMTP
5	База данных	Работа с PostgreSQL для хранения данных
</div>
📋 Предварительные требования
<div style="background: #f5f5f5; padding: 15px; border-radius: 5px; border-left: 4px solid #2ecc71;">
Go 1.19 или выше - язык программирования

PostgreSQL 12+ - система управления базами данных

Доступ к SMTP-серверу - для отправки email-уведомлений

Аккаунт в ЮКассе - для обработки платежей

</div>
⚙️ Настройка конфигурации
1. Создание файла конфигурации
Создайте файл config.yaml в папке config/ со следующим содержимым:

yaml
server:
  version: "Vitalis-api-v1.0.0" # Название и версия API
  port: ":8080"                 # порт на котором будет работать сервер

logger:
  level: "debug"                # уровень логирования: info, debug

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
2. Настройка переменных окружения
Создайте файл .env в папке cmd/ или экспортируйте переменные в среде выполнения:

env
# Настройки базы данных PostgreSQL
DB_USER=your_database_user
DB_PASSWORD=your_secure_password
DB_NAME=vitalis_db
DB_HOST=your_host
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
🗄️ Настройки базы данных
<table> <tr><th>Параметр</th><th>Описание</th><th>Пример</th></tr> <tr><td><code>DB_USER</code></td><td>Имя пользователя PostgreSQL</td><td><code>vitalis_user</code></td></tr> <tr><td><code>DB_PASSWORD</code></td><td>Пароль пользователя базы данных</td><td><code>secure_password_123</code></td></tr> <tr><td><code>DB_NAME</code></td><td>Название базы данных</td><td><code>vitalis_db</code></td></tr> <tr><td><code>DB_HOST</code></td><td>Хост базы данных</td><td><code>localhost</code></td></tr> <tr><td><code>DB_PORT</code></td><td>Порт для подключения к БД</td><td><code>5432</code></td></tr> </table>
📧 Настройки email-рассылок
<table> <tr><th>Параметр</th><th>Описание</th><th>Пример</th></tr> <tr><td><code>SMTP_HOST</code></td><td>Адрес SMTP-сервера для отправки писем</td><td><code>smtp.gmail.com</code></td></tr> <tr><td><code>SMTP_PORT</code></td><td>Порт SMTP-сервера</td><td><code>587</code></td></tr> <tr><td><code>SMTP_USER</code></td><td>Логин для авторизации на SMTP-сервере</td><td><code>your_email@gmail.com</code></td></tr> <tr><td><code>SMTP_PASSWORD</code></td><td>Пароль для авторизации на SMTP-сервере</td><td><code>app_specific_password</code></td></tr> </table>
💳 Настройки платежной системы
<table> <tr><th>Параметр</th><th>Описание</th><th>Где получить</th></tr> <tr><td><code>YOOKASSA_SHOP_ID</code></td><td>Идентификатор магазина в ЮКассе</td><td>Личный кабинет ЮКассы</td></tr> <tr><td><code>YOOKASSA_SECRET_KEY</code></td><td>Секретный ключ для доступа к API ЮКассы</td><td>Личный кабинет ЮКассы</td></tr> </table>
🌐 Общие настройки
<table> <tr><th>Параметр</th><th>Описание</th><th>Пример</th></tr> <tr><td><code>FRONTEND_URL</code></td><td>Базовый URL фронтенд-приложения</td><td><code>https://vitalis-life.ru</code></td></tr> <tr><td><code>MANAGER_EMAIL</code></td><td>Email адрес менеджера для уведомлений о заказах</td><td><code>manager@vitalis-life.ru</code></td></tr> </table>
🚀 Запуск приложения
<div style="background: #e8f4f8; padding: 15px; border-radius: 5px; border-left: 4px solid #3498db;">
1. Клонируйте репозиторий:
bash
git clone <repository-url>
cd vitalis-server
2. Установите зависимости:
bash
go mod download
3. Настройте конфигурацию (как описано выше)
4. Запустите приложение:
bash
go run cmd/main.go
Примечание: Приложение будет доступно по адресу: http://localhost:8080 (если локально запускаете)

</div>
📁 Структура проекта
<pre style="background: #2c3e50; color: #ecf0f1; padding: 15px; border-radius: 5px;"> vitalis-server/ ├── cmd/ # Точка входа приложения ├── config/ # Конфигурационные файлы ├── internal/ # Внутренние пакеты приложения │ ├── adapter/ # HTTP контроллеры │ ├── app/ # Бизнес-логика │ ├── domain/ # Работа с данными ├── pkg/ # Вспомогательные пакеты </pre>
📄 Лицензия
<div align="center" style="margin-top: 30px; padding: 20px; background: #f9f9f9; border-radius: 5px;">
© 2024 Vitalis Life. Все права защищены.

Этот проект является собственностью Vitalis Life. Распространение и использование без разрешения запрещено.

</div>
<div align="center">
Документация обновлена: 📅 {дата}

</div>
