# 🧮 Distributed Calculator
Распределённый калькулятор с поддержкой сложных выражений, асинхронной обработки и истории вычислений.
Считает даже ~(~2) + 3 и помнит, что 1000 - 7 = 993.

## ⚙️ Требования к окружению
Перед запуском проекта убедитесь, что у вас установлено:

`Docker` и `Docker Compose` (для запуска)
`Node.js 18+` и `npm` (для запуска фронтенда)
## 🚀 Запуск проекта
### 1. Подготовка окружения
```bash
# Создайте .env файл на основе примера
cp .env.example .env
```
### 2. Запуск бэкенда
```bash
# Собираем образ
docker-compose build --no-cache

# Запускаем стек
docker-compose up
```
### 3. Запуск фронтенда
```bash
cd my-calculator
npm install
npm run dev
```
## 🔨 Технологический стек
| Компонент | Технология |
| :---:  | :---:     |
| Сообщения | `Kafka` (3-нодовый кластер)|
| Кэш | `Redis` |
| БД | `PostgreSQL` |
| API | `gRPC` + `HTTP/JSON` |
| Авторизация | `JWT` |
| Фронтенд | `React` + `Vite` + `Tailwind CSS` |
| Сборка | `Docker Compose` |

## 🧠 Как это работает
1. Пользователь вводит выражение в веб-интерфейсе
2. Фронт отправляет запрос на `Gateway`
3. Gateway проверяет `JWT` и парсит выражение
4. Задача разбивается на шаги и отправляется в `Kafka`
5. Воркеры обрабатывают шаги, сохраняя промежуточные результаты в `Redis`
6. Финальный результат сохраняется в `PostgreSQL`
7. Пользователь получает результат или ошибку

## 📡 API Endpoints
| Метод | URL | Описание |
| :---: | :---: | :---: |
| `POST` | `/v1/calculate` | Запускает вычисление выражения |
| `POST` | `/v1/result`    | Возвращает результат по `task_id` |
| `POST` | `/v1/examples` | Возвращает историю вычислений пользователя |
| `POST` | `/v1/register` | Регистрация пользователя | 
| `POST` | `/v1/login`    | Авторизация и получение JWT |

### 💡 Примеры использования
✅ Пример: Регистрация пользователя </br>
Запрос
```bash
curl --location 'http://localhost:8080/v1/register' \
--header 'Content-Type: application/json' \
--data-raw '{
  "email": "user@example.com",
  "password": "mysecretpassword123"
}'
```
Ответ
```json
{
  "success": true,
  "error": ""
}
```
✅ Пример: Вход в систему </br>
Запрос
```bash
curl --location 'http://localhost:8080/v1/login' \
--header 'Content-Type: application/json' \
--data-raw '{
  "email": "user@example.com",
  "password": "mysecretpassword123"
}'
```
Ответ
```json
{
    "success": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZDk0ODE5NzEtMjMzMy00MzE2LWFkZDYtYWE1NmFlYzc1OTgxIiwiaXNzIjoiZGlzdHJpYnV0ZWRfY2FsY3VsYXRvciIsImV4cCI6MTc1Mzk4OTAxM30.pEd-1x3AfYreT4gzRWeo-oBcUzGdfDaBlUe31HIWddA",
    "userId": "d9481971-2333-4316-add6-aa56aec75981",
    "error": ""
}

```
✅ Пример: Вычисление выражения <br>
Запрос 
```bash
curl --location 'http://localhost:8080/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: ••••••' \
--data '{
    "expression": "1000-7"
}'
```
Ответ
```json
{
    "taskId": "c352230c-802e-4158-b528-5b2365481179"
}

```
✅ Пример: Получение результата </br>
Запрос 
```bash
curl --location 'http://localhost:8080/v1/result' \
--header 'Content-Type: application/json' \
--header 'Authorization: ••••••' \
--data '{
    "taskId": "c352230c-802e-4158-b528-5b2365481179"
}'

```
Ответ
```json
{
    "value": 993
}
```
✅ Пример: История вычислений </br>
Запрос
```bash
curl --location --request POST 'http://localhost:8080/v1/examples' \
--header 'Authorization: ••••••'
```
Ответ
```json
{
    "examples": [
        {
            "id": "2b19f0d8-5674-4de3-b0fc-1ad09db15572",
            "expression": "Hello, World!",
            "calculated": true,
            "createdAt": "2025-07-30T19:15:40Z",
            "error": "line is not a mathematical expression or contains an error"
        },
        {
            "id": "8f70e303-44e5-46de-b1da-13f460d455af",
            "expression": "7 / 0",
            "calculated": true,
            "createdAt": "2025-07-30T19:15:21Z",
            "error": "division by zero"
        },
        {
            "id": "d136bce9-06f6-470d-a088-bda9a2e132be",
            "expression": "~(~3) + 8 ^ 0",
            "calculated": true,
            "result": 4,
            "createdAt": "2025-07-30T19:15:00Z"
        },
        {
            "id": "a8f71354-fb00-4034-a9ea-9515caf7bd77",
            "expression": "~3 + 8",
            "calculated": true,
            "result": 5,
            "createdAt": "2025-07-30T19:14:28Z"
        },
        {
            "id": "c352230c-802e-4158-b528-5b2365481179",
            "expression": "1000-7",
            "calculated": true,
            "result": 993,
            "createdAt": "2025-07-30T19:11:37Z"
        }
    ]
}
```
## 🗂️ Структура проекта
```
distributed_calculator2/
├── cmd/
│   ├── main/        # Основной сервер (gRPC + REST)
│   └── worker/      # Воркер для обработки задач
├── internal/
│   ├── auth/        # JWT авторизация
│   ├── models/      # Модели данных
│   ├── repository/  # Репозитории (Postgres, Redis)
│   ├── service/     # Бизнес-логика
│   └── worker/      # Логика воркера
├── pkg/
│   ├── api/         # gRPC прото
│   ├── config/      # Конфигурация
│   ├── db/          # Подключение к БД
│   ├── logger/      # Логирование
│   ├── messaging/   # Kafka
│   └── valueprovider/ # Получение значений
├── migrations/      # Миграции БД
├── my-calculator/   # React фронтенд
├── docker-compose.yml
├── Dockerfile
└── .env.example
```
## 🗃️ Структура БД
### Таблица examples
| Поле | Тип | Описание |
| :---: | :---: | :---: |
|`id`|`TEXT`|Уникальный `ID` выражения
|`expression`|`TEXT`|Исходное выражение
|`response`|`TEXT`|Финальная переменная
|`user_id`|`TEXT`|`ID` пользователя
|`calculated`|`BOOLEAN`|Вычисление завершено
|`error`|`TEXT`|Ошибка (если есть)
|`created_at`|`TIMESTAMPTZ`|Время создания
|`updated_at`|`TIMESTAMPTZ`|Время обновления
### Таблица users
Поле|Тип|Описание
| :---: | :---: | :---: |
|`id`|`TEXT`|Уникальный ID пользователя
|`email`|`TEXT`|Email пользователя
|`password_hash`|`TEXT`|Хэш пароля
|`role`|`TEXT`|Роль (user/admin)
|`created_at`|`TIMESTAMPTZ`|Время регистрации
|`updated_at`|`TIMESTAMPTZ`|Время обновления

## 🧩 Особенности реализации
1. Унарный минус через ~
```
// ~5 превращается в (0-5)
// ~(~2) + 3 = 5
```
2. Обработка деления на ноль
```
5 / 0 → ошибка "division by zero"
```
3. Асинхронная обработка
* Выражение разбивается на шаги
* Каждый шаг отправляется в `Kafka`
* Воркеры обрабатывают шаги параллельно
* Результат собирается из промежуточных значений
4. Поддержка сложных выражений
```
~(~2) + 3 * (4 - 1) ^ 2
```
## 🖥️ Фронтенд
Фронтенд на `React` с темной темой (чёрный фон, фиолетовые акценты):

* Главная страница — описание проекта и технологии
* Калькулятор — ввод выражений и получение результатов
* История — просмотр предыдущих вычислений
* Авторизация — регистрация и вход
* Доступен на [`http://localhost:5173`](http://localhost:5173) после запуска 
```bash
npm run dev
```

### 🛠️ Как считается выражение
1. Пользователь вводит выражение: `~(~2) + 3`
2. Система парсит его в обратную польскую нотацию:
```
2 ~ 2 ~ 3 + → 2 (0 - 2) (0 - 3) + 
```
#### Разбивает на шаги:
* Шаг 1: `~2 = -2`
* Шаг 2: `~(-2) = 2`
* Шаг 3: `2 + 3 = 5`
* Каждый шаг отправляется в `Kafka`
* Воркеры обрабатывают шаги и сохраняют результаты в `Redis`
* Финальный результат сохраняется в `PostgreSQL`
## 📊 Мониторинг воркеров
В фронтенде есть раздел "`Workers`", где отображается:

* Состояние воркеров (онлайн/оффлайн)
* Количество обработанных задач
* Текущая нагрузка
* Время последней активности
## 🔐 Безопасность
* Все запросы требуют `JWT`-авторизации (кроме регистрации и входа)
* Пароли хранятся в хэшированном виде (`bcrypt`)
* Валидация входных данных на всех этапах
* Нет использования `eval()` — безопасный парсинг выражений
## 📁 Docker Compose
Проект использует мощный `Docker Compose` с:

* 3-нодовым кластером `Kafka` (`KRaft`)
* `Redis` для кэширования
* `PostgreSQL` для хранения истории
* Автоматическими миграциями
* Созданием топиков `Kafka` при старте