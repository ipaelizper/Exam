package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./exam.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// สร้างตาราง
	createTables()
}

func createTables() {
	// ตารางคำถาม
	questionTable := `
	CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL
	);`

	// ตารางคำตอบ
	answerTable := `
	CREATE TABLE IF NOT EXISTS answers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		question_id INTEGER,
		text TEXT NOT NULL,
		FOREIGN KEY (question_id) REFERENCES questions (id) ON DELETE CASCADE
	);`

	_, err := db.Exec(questionTable)
	if err != nil {
		log.Fatal("Failed to create questions table:", err)
	}

	_, err = db.Exec(answerTable)
	if err != nil {
		log.Fatal("Failed to create answers table:", err)
	}
}
