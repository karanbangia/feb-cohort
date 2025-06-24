# get isolation level
SELECT @@transaction_isolation;
#default in mysql is REPEATABLE READ
SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ;

START TRANSACTION;

SELECT * FROM users WHERE id=1 FOR UPDATE;

UPDATE users SET name='john' WHERE id=1;

SELECT * FROM users WHERE id=1;

COMMIT;


SELECT * FROM users WHERE id=1;