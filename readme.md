## velesbok
Учебное веб-приложение, в перспективе с заметками
```md
velesbook/
│── cmd/               # Точка входа приложения
│   ├── server/        # Логика запуска HTTP-сервера
│   │   ├── server.go  # Инициализация маршрутов и серверных зависимостей
│   │   ├── routes.go  # Определение маршрутов
│── config/            # Конфигурационные файлы
│   ├── config.go      # Загрузка переменных окружения (env)
│── internal/          # Основная логика приложения
│   ├── auth/          # Логика аутентификации (JWT, bcrypt)
│   ├── user/          # Работа с пользователями
│   ├── page/          # Логика работы со страницами
│   ├── database/      # Подключение к БД
│── migrations/        # SQL-файлы для миграций
│── pkg/               # Вспомогательные утилиты
│── main.go            # Главная точка входа
│── go.mod             # Файл модулей
│── Dockerfile         # Контейнеризация
```


```json
{
  "endpoints": {
    "POST /auth/register": {
      "body": {
        "email": "user@example.com",
        "password": "securepassword"
      }
    },
    "POST /auth/login": {
      "body": {
        "email": "user@example.com",
        "password": "securepassword"
      }
    },
    "AuthMiddleware": {
      "headers": {
        "Authorization": "Bearer <token>"
      }
    },
    "GET /users": {
      "headers": {
        "Authorization": "Bearer <token>"
      }
    },
    "GET /pages": {
      "headers": {
        "Authorization": "Bearer <token>"
      }
    },
    "POST /pages/create": {
      "headers": {
        "Authorization": "Bearer <token>"
      },
      "body": {
        "title": "Название страницы",
        "content": "Текст страницы",
        "user_id": 1,
        "parent_id": null
      }
    }
  }
}
```