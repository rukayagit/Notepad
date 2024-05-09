package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/gorilla/mux"
	"C:\\Users\\Рукайя\\Desktop\\Notepad\\pkg\\note"
)

func main() {
	note.InitDB("user=postgres dbname=note_db sslmode=disable")

	r := mux.NewRouter()

	r.HandleFunc("/notes", CreateNoteHandler).Methods("POST")
	r.HandleFunc("/notes/{id}", GetNoteHandler).Methods("GET")
	r.HandleFunc("/notes/{id}", UpdateNoteHandler).Methods("PUT")
	r.HandleFunc("/notes/{id}", DeleteNoteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка при чтении данных запроса", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	newNote := note.Note{
		Title:   title,
		Content: content,
	}

	id, err := note.CreateNote(newNote)
	if err != nil {
		http.Error(w, "Ошибка при создании заметки", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Заметка успешно создана. ID: %d", id)
}

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Не указан ID заметки", http.StatusBadRequest)
		return
	}

	noteID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID заметки должен быть целым числом", http.StatusBadRequest)
		return
	}

	note, err := note.GetNoteById(noteID)
	if err != nil {
		http.Error(w, "Ошибка при получении заметки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка при чтении данных запроса", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Не указан ID заметки", http.StatusBadRequest)
		return
	}

	noteID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID заметки должен быть целым числом", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	updatedNote := note.Note{
		ID:      noteID,
		Title:   title,
		Content: content,
	}

	err = note.UpdateNote(updatedNote)
	if err != nil {
		http.Error(w, "Ошибка при обновлении заметки", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Заметка успешно обновлена")
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Не указан ID заметки", http.StatusBadRequest)
		return
	}

	noteID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID заметки должен быть целым числом", http.StatusBadRequest)
		return
	}

	err = note.DeleteNote(noteID)
	if err != nil {
		http.Error(w, "Ошибка при удалении заметки", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Заметка успешно удалена")
}
