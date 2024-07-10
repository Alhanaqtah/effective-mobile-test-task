# Time Tracker API

Этот проект представляет собой API для отслеживания задач. Он реализован на языке Go с использованием фреймворка `go-chi` и генерирует документацию Swagger с помощью библиотеки `swaggo/swag`.

## Установка

1. Клонируйте репозиторий:
    ```sh
    git clone https://github.com/yourusername/time-tracker.git
    cd time-tracker
    ```

2. Создайте файл `.env` в корневом каталоге проекта со следущими переменными:
    ```env
    ENV=local # local/env/prod

    POSTGRES_USER=
    POSTGRES_PASSWORD=
    POSTGRES_HOST=
    POSTGRES_PORT=
    DATABASE_NAME=time_tracker_db

    SERVER_HOST=
    SERVER_PORT=
    SERVER_TIMEOUT=

    EXTERNAL_API_URL=
    ```

3. Установите зависимости:
    ```sh
    go mod tidy
    ```

## Запуск

1. Запустите PostgreSQL и создайте базу данных `time_tracker_db`.
   
2. Запустите сервер:
    ```sh
    make run
    ```


## Документация API

Документация Swagger доступна по адресу `/docs`. Вы можете использовать её для тестирования и ознакомления с доступными конечными точками API.
