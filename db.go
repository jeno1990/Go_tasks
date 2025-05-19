package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMysqlStorage(cfg mysql.Config) (*sql.DB, error) {
	dsn := cfg.FormatDSN() // this is the driver to sql
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func addTask(db *sql.DB, task Task) (sql.Result, error) {
	query := "INSERT INTO tasks (name, description, status, created_at) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, task.Name, task.Description, task.Status, task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			status BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

func getTasks(db *sql.DB) ([]Task, error) {
	query := "SELECT id, name, description, status, created_at FROM tasks"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []Task{}
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func getTask(db *sql.DB, id int) (Task, error) {
	query := "SELECT id, name, description, status, created_at FROM tasks WHERE id = ?"
	row := db.QueryRow(query, id)
	var task Task
	err := row.Scan(&task.ID, &task.Name, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func deleteTask(db *sql.DB, id int) (sql.Result, error) {
	query := "DELETE FROM tasks WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func completeTask(db *sql.DB, id int) (sql.Result, error) {
	query := "UPDATE tasks SET status = true WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
