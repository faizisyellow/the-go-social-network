ALTER TABLE users ADD COLUMN role_id INT DEFAULT 1;

ALTER TABLE users ADD CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id);

UPDATE users SET role_id = (SELECT id FROM roles WHERE name = 'user');

ALTER TABLE users MODIFY COLUMN role_id INT NOT NULL;
    