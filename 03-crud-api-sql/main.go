package main

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    _ "github.com/mattn/go-sqlite3"
    _ "net/http/pprof"
)

type User struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
    Email   string  `json:"email"`
}


func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/users", listUsersHandler)
    mux.HandleFunc("POST /users", createUserHandler)
    mux.HandleFunc("/cpu", cpuIntensiveHandler)
    go http.ListenAndServe(":3000", mux)
    http.ListenAndServe(":6060", nil)
}

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("sqlite3", "users.db")
    if err != nil { 
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    rows, err := db.Query("Select * from users")
    if err != nil { 
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    users := []User{}
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }
    if err := json.NewEncoder(w).Encode(users); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("sqlite3", "users.db")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var u User

    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if _, err := db.Exec("INSERT INTO users (id, name, email) VALUES (?, ?, ?)", u.ID, u.Name, u.Email); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}

func cpuIntensiveHandler(w http.ResponseWriter, r *http.Request) {
    result := fibonacci(60)
    w.Write([]byte(strconv.Itoa(result)))
}
