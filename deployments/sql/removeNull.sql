UPDATE user SET salt = "" WHERE salt IS NULL;
UPDATE user SET bio = "" WHERE bio IS NULL;
UPDATE tide SET about = "" WHERE about IS NULL;