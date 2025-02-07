CREATE TABLE pages (
    id SERIAL PRIMARY KEY,                     -- Автоинкремент для id
    title VARCHAR(255) NOT NULL,                -- Заголовок страницы
    content TEXT,                              -- Содержание страницы
    user_id BIGINT NOT NULL,                   -- Внешний ключ на пользователя
    parent_id BIGINT,                          -- Внешний ключ на родительскую страницу
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время создания страницы
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время последнего обновления страницы
    FOREIGN KEY (parent_id) REFERENCES pages(id) ON DELETE CASCADE, -- Родитель ссылается на id этой же таблицы
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE   -- Связь с пользователями
);