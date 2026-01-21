package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbHost     = "host.docker.internal"
	dbPort     = 3306
	dbUser     = "root"
	dbPassword = "12345"
	dbName     = "sten_sax_pase"
)

func initDatabase() *sql.DB {

	// Bygga vår DSN
	serverDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		dbUser,dbPassword, dbHost, dbPort,
	)

	serverDb, err := sql.Open("mysql", serverDsn)

	if err != nil {
		log.Fatal("Kunde inte skapa serverDB: ", err)
	}

	if err := serverDb.Ping(); err != nil {
		log.Fatal("Kunde inte ansluta till mysql-servern: ", err)
	}

	_, err = serverDb.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName,
	))

	if err != nil {
		log.Fatal("Kunde inte skapa databasen", err)
	}	

	_ = serverDb.Close()


	appDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		dbUser,dbPassword, dbHost, dbPort, dbName,
	)

	db, err := sql.Open("mysql", appDSN)
	if err != nil {
		log.Fatal("Kunde inte skapa db")
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Kunde inte ansluta till databasen: ", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS rounds (
			id INT AUTO_INCREMENT PRIMARY KEY,
			played_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			player_choice VARCHAR(10) NOT NULL,
			computer_choice VARCHAR(10) NOT NULL,
			result VARCHAR(10) NOT NULL
		)
	`)

	if err != nil {
		log.Fatal("Kunde inte skapa tabellen rounds", err)
	}

	return db

}




// Spara rounds

func insertRound(db *sql.DB, playerChoice string, computerChoice string, result string) error {
	// Göra en insert

	_, err := db.Exec(
		"INSERT INTO rounds(player_choice, computer_choice, result) VALUES(?,?,?)",
		playerChoice, computerChoice, result,
	)
	return err
}


type Stats struct {
	Total int
	Wins int
	WinPct int
}

// Hämta stats

func getStats(db *sql.DB)(Stats, error) {
	var total int
	var wins int

	err := db.QueryRow("SELECT COUNT(*) FROM rounds").Scan(&total)
	if err != nil {
		return Stats{}, err
	}

	err = db.QueryRow("SELECT COUNT(*) FROM rounds WHERE result = 'vinst'").Scan(&wins)
	if err != nil {
		return Stats{}, err
	}

	winPct := 0
	
	if total > 0 {
		winPct = int(float64(wins) / float64(total) * 100.0)
	}

	return Stats{Total: total, Wins: wins, WinPct: winPct}, nil

}


