CREATE TABLE IF NOT EXISTS user_invitations(
    token VARBINARY(72) NOT NULL PRIMARY KEY,
    user_id INT NOT NULL
);