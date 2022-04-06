package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// ImportCsvPG : copy csv file into database.
func ImportCsvPgSales(srcFile string) {
	c := ParseConfig().PGConn
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	r1, err := db.Exec("DELETE FROM sales_raw")
	if err != nil {
		log.Fatal("Failed to delete data from sales_raw: ", err)
	}
	affected1, _ := r1.RowsAffected()
	fmt.Printf("Deleted  : %d rows\n", affected1)

	r2, err := db.Exec(fmt.Sprintf("COPY sales_raw FROM '%s' DELIMITER ',' CSV HEADER", srcFile))
	if err != nil {
		log.Fatal("Failed to import csv file: ", err)
	}
	affected2, _ := r2.RowsAffected()
	fmt.Printf("Inserted to sales_raw table: %d rows\n", affected2)

	sqlStr := `
    INSERT INTO sales
    SELECT salesdate, agencytype, fop, fopdesc, domintl, salesrefund, ccy, sum(amount) amount, sum(krwamt) krwamt
    FROM sales_raw GROUP BY 1,2,3,4,5,6,7
    ON CONFLICT ON CONSTRAINT sales_pkey 
    DO UPDATE SET amount = EXCLUDED.amount, krwamt = EXCLUDED.krwamt;`
	affected3, err := db.Exec(sqlStr)
	if err != nil {
		log.Fatal("Failed to insert summarized data into sales table: ", err)
	}
	fmt.Printf("Inserted to sales table: %v rows\n", affected3)

	_, err = db.Exec("VACUUM ANALYZE")
	if err != nil {
		log.Fatal("Failed to VACUUM: ", err)
	}
	fmt.Printf("Data (sales) import finished succeffully\n")
}

func ImportCsvPgTaxYr(srcFile string) {
	c := ParseConfig().PGConn
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	r1, err := db.Exec("DELETE FROM taxyr_raw")
	if err != nil {
		log.Fatal("Failed to delete data from taxyr_raw: ", err)
	}
	affected1, _ := r1.RowsAffected()
	fmt.Printf("Deleted  : %d rows\n", affected1)

	r2, err := db.Exec(fmt.Sprintf("COPY taxyr_raw FROM '%s' DELIMITER ',' CSV HEADER", srcFile))
	if err != nil {
		log.Fatal("Failed to import csv file: ", err)
	}
	affected2, _ := r2.RowsAffected()
	fmt.Printf("Inserted to taxyr_raw table: %d rows\n", affected2)

	sqlStr := `
    INSERT INTO taxyr
	SELECT salesdate, taxyr, domintl, ccy, 
	       sum(salesamt) salesamt, sum(refundamt) refundamt, sum(reissueamt) reissueamt
	FROM taxyr_raw
	GROUP BY salesdate, taxyr, domintl, ccy
	ORDER BY salesdate, taxyr, domintl, ccy
	ON CONFLICT ON CONSTRAINT taxyr_pkey 
	DO UPDATE set salesamt = excluded.salesamt, refundamt = excluded.refundamt, reissueamt = excluded.reissueamt;`
	affected3, err := db.Exec(sqlStr)
	if err != nil {
		log.Fatal("Failed to insert summarized data into taxyr table: ", err)
	}
	fmt.Printf("Inserted to staxyr table: %v rows\n", affected3)

	_, err = db.Exec("VACUUM ANALYZE")
	if err != nil {
		log.Fatal("Failed to VACUUM: ", err)
	}
	fmt.Printf("Data (taxyr) import finished succeffully\n")
}