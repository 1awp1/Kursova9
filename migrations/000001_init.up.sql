CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login varchar(255) not null,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    pass_hash TEXT NOT NULL,
    token TEXT,
    role_id UUID,
    phone_number VARCHAR(255),
    email VARCHAR(255),
    status boolean,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL 
);

INSERT INTO roles(role_name) VALUES ('employee'), ('manager'), ('admin');

WITH manager_role AS (
    SELECT id FROM roles WHERE role_name = 'manager'
),
admin_role AS (
    SELECT id FROM roles WHERE role_name = 'admin'
)

INSERT INTO users (login, pass_hash, role_id, phone_number, email)
VALUES
    ('manag', '$2a$10$eYnJgFmQwIY5Jja5uR0.4ut3xLlL6yq3IjxIfqDwRLMM7VFxi9zT6', (SELECT id FROM manager_role), '89228990747', 'deutchwar@gmail.com'),
    ('admi', '$2a$10$56x4DjRzGq1ersvqKuXgfeXdlczik0MzP0lXt9NvalpW20O1QjdBW', (SELECT id FROM admin_role), '89228990747', 'deutchwar@gmail.com');
