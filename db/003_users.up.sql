ALTER TABLE users ADD CONSTRAINT unique_user_id UNIQUE (user_id);
SELECT * FROM users;

DELETE FROM users WHERE user_id = 873923391;
