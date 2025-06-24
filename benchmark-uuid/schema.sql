CREATE TABLE UUID_BENCHMARK
(
    id   VARCHAR(256) NOT NULL,
    PRIMARY KEY (id)
);

CREATE  TABLE NUMBER_BENCHMARK
(
    id INTEGER NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (id)
);

ALTER TABLE NUMBER_BENCHMARK ADD COLUMN stub VARCHAR(255);


SELECT  COUNT(*) FROM UUID_BENCHMARK UNION
SELECT  COUNT(*) FROM NUMBER_BENCHMARK;

select
    database_name,
    table_name,
    index_name,
    round((stat_value*@@innodb_page_size)/1024/1024, 2) SizeMB,
    round(((100/(SELECT INDEX_LENGTH FROM INFORMATION_SCHEMA.TABLES t WHERE t.TABLE_NAME = iis.table_name and t.TABLE_SCHEMA = iis.database_name))*(stat_value*@@innodb_page_size)), 2) `Percentage`
from mysql.innodb_index_stats iis
where stat_name='size'
  and table_name = 'NUMBER_BENCHMARK'
  and database_name = 'mysql';

select
    database_name,
    table_name,
    index_name,
    round((stat_value*@@innodb_page_size)/1024/1024, 2) SizeMB,
    round(((100/(SELECT INDEX_LENGTH FROM INFORMATION_SCHEMA.TABLES t WHERE t.TABLE_NAME = iis.table_name and t.TABLE_SCHEMA = iis.database_name))*(stat_value*@@innodb_page_size)), 2) `Percentage`
from mysql.innodb_index_stats iis
where stat_name='size'
  and table_name = 'UUID_BENCHMARK'
  and database_name = 'mysql';