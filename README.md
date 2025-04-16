# Не Стать Тобой (Ne Stat Toboy)

Сайт-визитка для короткометражного фильма "Не Стать Тобой".

## Технологии

- Go 1.24
- Echo v4 (веб-фреймворк)
- templ (для генерации HTML)
- htmx (для интерактивности на фронтенде)

## Структура проекта

```
ne-stat-toboy/
├── cmd/
│   └── server/        # Main server code
├── internal/
│   ├── handler/       # HTTP handlers
│   └── model/         # Data models
└── web/
    ├── static/        # Static files
    │   ├── css/       # Stylesheets
    │   └── img/       # Images
    └── template/      # templ templates
```

## Установка и запуск

### Предварительные требования

- Go 1.24 или выше
- templ 0.3.857 или выше

### Установка templ

```bash
go install github.com/a-h/templ/cmd/templ@latest
```

### Запуск проекта

1. Клонируйте репозиторий:

```bash
git clone https://github.com/lexfrei/ne-stat-toboy.git
cd ne-stat-toboy
```

2. Генерация templ-шаблонов:

```bash
templ generate
```

3. Установка зависимостей:

```bash
go mod download
```

4. Запустите сервер:

```bash
go run cmd/server/main.go
```

5. Откройте браузер и перейдите по адресу `http://localhost:8080`

### Сборка исполняемого файла

```bash
go build -o ne-stat-toboy ./cmd/server
```

### Запуск через Docker

```bash
docker build -t ne-stat-toboy .
docker run -p 8080:8080 ne-stat-toboy
```

## Структура сайта

- **Главная** — Основная информация о фильме
- **О фильме** — Синопсис и художественная концепция
- **Команда** — Информация о команде и актерах
- **Локации** — Описание локаций съемок
- **Контакты** — Контактная информация и форма обратной связи