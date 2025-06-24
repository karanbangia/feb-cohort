CREATE DATABASE DB_1;
CREATE DATABASE DB_2;

CREATE TABLE DB_1.USER (
    ID INT PRIMARY KEY,
    NAME VARCHAR(100)
);

CREATE TABLE DB_1.USER_PROFILE (
    ID INT PRIMARY KEY,
    USER_ID INT,
    PROFILE VARCHAR(100),
    FOREIGN KEY (USER_ID) REFERENCES DB_1.USER(ID)
);

CREATE TABLE DB_1.POST (
    ID INT PRIMARY KEY,
    USER_ID INT,
    TITLE VARCHAR(100),
    FOREIGN KEY (USER_ID) REFERENCES DB_1.USER(ID)
);

CREATE TABLE DB_2.USER (
                           ID INT PRIMARY KEY,
                           NAME VARCHAR(100)
);

CREATE TABLE DB_2.USER_PROFILE (
                                   ID INT PRIMARY KEY,
                                   USER_ID INT,
                                   PROFILE VARCHAR(100),
                                   FOREIGN KEY (USER_ID) REFERENCES DB_2.USER(ID)
);

CREATE TABLE DB_2.POST (
                           ID INT PRIMARY KEY,
                           USER_ID INT,
                           TITLE VARCHAR(100),
                           FOREIGN KEY (USER_ID) REFERENCES DB_2.USER(ID)
);

ALTER TABLE DB_2.USER ADD COLUMN EMAIL VARCHAR(100);


mysqldump -u dbeaver -p DB_1 > source_db_dump.sql
mysql -u dbeaver -p DB_2 < source_db_dump.sql

TRUNCATE TABLE DB_2.USER;

mysqldump -u dbeaver -p --no-create-info DB_1 > data_dump.sql
