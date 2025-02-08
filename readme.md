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
    "POST /v1/auth/register": {
      "body": {
        "email": "user@example.com",
        "password": "securepassword"
      }
    },
    "POST /v1/auth/login": {
      "body": {
        "email": "user@example.com",
        "password": "securepassword"
      }
    },
    "GET /v1/api/users": {
      "headers": {
        "Authorization": "Bearer <token>"
      }
    },
    "GET /v1/api/pages": {
      "headers": {
        "Authorization": "Bearer <token>"
      }
    },
    "POST /v1/api/pages/create": {
      "headers": {
        "Authorization": "Bearer <token>"
      },
      "body": {
        "title": "Название страницы",
        "content": "Текст страницы",
        "parent_id": null
      }
    }
  }
}
```