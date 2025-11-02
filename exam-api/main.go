package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Answer struct {
	ID   int    `json:"id,omitempty"`
	Text string `json:"text"`
}

type Question struct {
	ID      int      `json:"id"`
	Number  int      `json:"number"`
	Content string   `json:"content"`
	Answers []Answer `json:"answers"`
}

var mu sync.Mutex

// getข้อสอบทั้งหมด
func getQuestions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	rows, err := db.Query(`
		SELECT q.id, q.content, a.id, a.text
		FROM questions q
		LEFT JOIN answers a ON q.id = a.question_id
		ORDER BY q.id, a.id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	questionMap := make(map[int]*Question)
	var order int

	for rows.Next() {
		var qID int
		var content string
		var aID sql.NullInt64
		var aText sql.NullString

		err := rows.Scan(&qID, &content, &aID, &aText)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// สร้างคำถามใหม่ ถ้ายังไม่มี
		if _, exists := questionMap[qID]; !exists {
			order++
			q := &Question{
				ID:      qID,
				Number:  order,
				Content: content,
				Answers: []Answer{},
			}
			questionMap[qID] = q
		}

		// เพิ่มคำตอบ
		if aID.Valid && aText.Valid {
			questionMap[qID].Answers = append(questionMap[qID].Answers, Answer{
				ID:   int(aID.Int64),
				Text: aText.String,
			})
		}
	}

	var questions []Question
	for _, q := range questionMap {
		questions = append(questions, *q)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// เพิ่มข้อสอบ
func addQuestion(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var input struct {
		Content string   `json:"content"`
		Answers []Answer `json:"answers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//  transaction
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// เพิ่มคำถาม
	res, err := tx.Exec("INSERT INTO questions (content) VALUES (?)", input.Content)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	questionID, _ := res.LastInsertId()

	// เพิ่มคำตอบ
	stmt, err := tx.Prepare("INSERT INTO answers (question_id, text) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, ans := range input.Answers {
		_, err := stmt.Exec(questionID, ans.Text)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ส่งกลับคำถามที่เพิ่ม
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": questionID})
}

// ลบข้อสอบ
func deleteQuestion(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM questions WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// เริ่มดาต้าเบส
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/questions", getQuestions).Methods("GET")
	r.HandleFunc("/questions", addQuestion).Methods("POST")
	r.HandleFunc("/questions/{id}", deleteQuestion).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
	})

	handler := c.Handler(r)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Database: exam.db")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
