-- Подключаем расширение для UUID (если еще не подключено)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    login VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    pass_hash TEXT NOT NULL,
    is_online BOOLEAN,
    role_id UUID, 
    phone_number VARCHAR(255),
    email VARCHAR(255),
    status BOOLEAN,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL
);

-- Вставляем роли
INSERT INTO roles(role_name) 
VALUES 
    ('user'),
    ('manager'),
    ('admin')
ON CONFLICT (role_name) DO NOTHING;

-- Добавляем пользователей
DO $$
DECLARE
    manager_id UUID;
    admin_id UUID;
BEGIN
    SELECT id INTO manager_id FROM roles WHERE role_name = 'manager';
    SELECT id INTO admin_id FROM roles WHERE role_name = 'admin';

    INSERT INTO users (login, first_name, last_name, pass_hash, role_id, phone_number, email)
    VALUES
        ('manag', 'John', 'Doe', '$2a$10$eYnJgFmQwIY5Jja5uR0.4ut3xLlL6yq3IjxIfqDwRLMM7VFxi9zT6', manager_id, '89228990747', 'deutchwar@gmail.com'),
        ('admi', 'Admin', 'User', '$2a$10$56x4DjRzGq1ersvqKuXgfeXdlczik0MzP0lXt9NvalpW20O1QjdBW', admin_id, '89228990747', 'deutchwar@gmail.com');
END $$;
