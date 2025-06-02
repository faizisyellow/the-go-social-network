CREATE TABLE IF NOT EXISTS roles(
    id  INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description VARCHAR(255)
);

INSERT INTO roles (name,level,description) VALUES
('user',1,'A user can create posts and comments'),
('moderator',2,'A moderator can update other users posts'),
('admin',3,'An admin can update and delete other users posts');