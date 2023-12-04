package main

import (
	"Web-Api/models"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func signupHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponseSignup{Error: "Invalid request format"})
		return
	}

	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", newUser.Name, newUser.Email, newUser.Password)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponseSignup{Error: "Error storing user in the database"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.SuccessResponseSignup{Message: "User created successfully"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid request format"})
		return
	}

	id, err := authenticateUser(requestData["email"].(string), requestData["password"].(string))
	if err != nil {

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid email or password"})
		return
	}

	w.WriteHeader(http.StatusOK)
	successResponse := models.SuccessResponse{SID: id}
	json.NewEncoder(w).Encode(successResponse)
}

type NoteCreationResponse struct {
	ID uint32 `json:"id"`
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var noteData map[string]string
	err := json.NewDecoder(r.Body).Decode(&noteData)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid request format"})
		return
	}

	sid := noteData["sid"]

	userID, err := authenticateUserByID(sid)
	if err != nil {

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	result, err := db.Exec("INSERT INTO notes (user_id, note) VALUES (?, ?)", userID, noteData["note"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Error storing note in the database"})
		return
	}

	noteID, _ := result.LastInsertId()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NoteCreationResponse{ID: uint32(noteID)})
}

func authenticateUserByID(sid string) (int, error) {

	var id int
	err := db.QueryRow("SELECT id FROM users WHERE id = ?", sid).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func authenticateUser(email, password string) (int, error) {

	var id int
	err := db.QueryRow("SELECT id FROM users WHERE email = ? AND password = ?", email, password).Scan(&id)
	if err != nil {
		return 0,
			fmt.Errorf("invalid credentials")
	}
	return id, nil
}

type NotesResponse struct {
	Notes []NoteInfo `json:"notes"`
}

type NoteInfo struct {
	ID   uint32 `json:"id"`
	Note string `json:"note"`
}

func listNotesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid request format"})
		return
	}

	sid := requestData["sid"]

	userID, err := authenticateUserByID(sid)
	if err != nil {

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	notes, err := getNotesByUserID(userID)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Error retrieving notes from the database"})
		return
	}

	w.WriteHeader(http.StatusOK)
	notesResponse := NotesResponse{Notes: notes}
	json.NewEncoder(w).Encode(notesResponse)
}

func getNotesByUserID(userID int) ([]NoteInfo, error) {
	rows, err := db.Query("SELECT id, note FROM notes WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []NoteInfo
	for rows.Next() {
		var note NoteInfo
		err := rows.Scan(&note.ID, &note.Note)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {

	var noteData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&noteData)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid request format"})
		return
	}

	sid, ok := noteData["sid"].(string)
	if !ok {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	noteIDFloat, ok := noteData["id"].(float64)
	if !ok {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid note ID format"})
		return
	}
	noteID := uint32(noteIDFloat)

	userID, err := authenticateUserByID(sid)
	if err != nil {

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notes WHERE id = ? AND user_id = ?", noteID, userID).Scan(&count)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Error checking note existence in the database"})
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Note ID not found for the user"})
		return
	}

	_, err = db.Exec("DELETE FROM notes WHERE id = ? AND user_id = ?", noteID, userID)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Error deleting note from the database"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.SuccessResponse{SID: userID})
}

func main() {

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", signupHandler)

	http.HandleFunc("/notes", noteHandler)

	fmt.Println("Server is listening on :8083...")
	err := http.ListenAndServe(":8083", nil)
	if err != nil {

		fmt.Println("Error starting server:", err)
	}
}

func noteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listNotesHandler(w, r)
	case http.MethodPost:

		createNoteHandler(w, r)
	case http.MethodDelete:
		deleteNoteHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
