# sales-info
Utility program for weekly sales report.

## Sub programs
- main.go
- parse-config.go 
  . parsing configuration from config.json
- parse-flag.fo   
  . parsing command line interface options
- convert-data.go
  . convert vectis download data --> import-data.csv
- query-postgres.go 
  . query sales results, export query results (PostgreSQL backend)
- query-sqlite.go 
  . query sales results, export query results (Sqlite3 backend)

## Files
- config.json
  . this file has default configuration setting
- import.sql
  . sqlite3 batch file for data import

## Config.json
{
    "root_dir": "D:/Projects/sales-info",
    "data": {
        "dir_name": "D:/Projects/sales-info/data",
        "source_file": "vectis_sales.csv",
        "import_file": "import_data.csv",
        "export_file": "sales_results.csv"
    },
    "database": "Sqlite",
    "pgconn": {
        "host": "localhost",
        "port": 5432,
        "user": "********",
        "password": "*******",
        "dbname": "selabd"
    },
    "sqlite_db": "D:/Projects/sales-info/data/selabd.db"
}

## Import.sql
.separator ","
.import D:/Projects/sales-info/data/import_data.csv sales --skip 1

- file name in import.sql should be equal to import_file in config.json 

## How to use

Usage: sales-info [-rpt table|raw] query  | [-in  filename -out filename] convert |
                  [-src filename ] import | [-dst filename] export

