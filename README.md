# sales-info
Utility program for weekly sales report.

## Sub programs
- main.go
- parse-config.go 
  . parsing configuration from config.json
- parse-flag.fo   
  . parsing command line interface options
- download-data.go
  . download vectis sales data using sales_download.exe
- convert-data.go
  . convert vectis download data --> import-data.csv
- query-postgres.go 
  . query sales results, export query results (PostgreSQL backend)
- query-sqlite.go 
  . query sales results, export query results (Sqlite3 backend)

## Files
Following files should exist in the same folder with sales-info.exe

- config.json
  . this file has default configuration setting
- import.sql
  . sqlite3 batch file for data import
- sqlite3.exe
- sales_download.exe

## Config.json
{
    "data": {
        "dir_name": "./data",
        "source_file": "VectisReport.csv",
        "import_file": "ImportData.csv",
        "export_file": "SalesResults.csv"
    },
    "database": "Sqlite",
    "pgconn": {
        "host": "localhost",
        "port": 5432,
        "user": "*********",
        "password": "*********",
        "dbname": "selabd"
    },
    "sqlite_db": "./data/selabd.db"
}

## Import.sql
.separator ","
.import ./data/Import_Data.csv sales --skip 1

- file name in import.sql should be equal to import_file in config.json 

## How to use

Usage: sales-info -from yyyy-mm-dd -to yyyy-mm-dd download |
       [-in filename -out filename] convert | [-src filename] import |
       [-rpt table|raw] query | [-dst filename] export | all

