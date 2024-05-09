package note

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Подключено к БД")
}

func CreateNote(note Note) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO notes (title, content, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", note.Title, note.Content, time.Now(), time.Now()).Scan(&id)
	if err!= nil {
		return 0, err
	}
	return id, nil
}

func GetNoteById(id int) (Note, error) {
	var note Note
	row := db.QueryRow("SELECT id, title, content, created_at, updated_at FROM notes WHERE id = $1", id)
	err := row.Scan(&note.Id, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err!= nil {
		return Note{}, err
	}
	return note, nil
}

func UpdateNote(note Note) error  {
	_, err := db.Exec("UPDATE notes SET title = $1, content = $2, updated_at = $3 WHERE id = $4", note.Title, note.Content, time.Now(), note.Id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNote(id int) error {
	_, err := db.Exec("DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}