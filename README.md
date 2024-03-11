# Importer
### Usage
```bash
make
./bin/importer
```
### Description
This program reads a *dirty* CSV file that parses to find errors in lines and then writes the correct lines to a new file.

The program has two options for the user to choose from:
1. Cleaner (clean dirty CSV file)
2. Importer (import clean CSV file)

You will need a .env file with the following variables:
```bash
MYSQL_USER=[username]
MYSQL_PASSWORD=[password]
MYSQL_DB=[database name]
CLEAN_DATA_FILE=data/clean/clean_data.csv
CLEAN_DIR=data/clean
DIRTY_DATA_FILE=data/dirty/dirty_data.csv
DIRTY_DIR=data/dirty #if you want to change the directory
ERROR_LOG_FILE=logs/errors.log
BAD_DATA_FILE=logs/bad_data.log
```


#### Database
The following is the database schema:
```sql
CREATE TABLE IF NOT EXISTS manufacturers (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS device_types (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS serial_numbers (
    auto_id INTEGER PRIMARY KEY,
    manufacturer_id INTEGER NOT NULL,
    device_type_id INTEGER NOT NULL,
    serial_number VARCHAR(68) NOT NULL,
    FOREIGN KEY (manufacturer_id) REFERENCES manufacturers(id),
    FOREIGN KEY (device_type_id) REFERENCES device_types(id)
);
```
### External Dependencies
- [sql-driver](https://github.com/go-sql-driver/mysql)
- [godotenv](https://github.com/joho/godotenv)
