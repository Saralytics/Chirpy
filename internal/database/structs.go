package database

import "sync"

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}
type DB struct {
	path string
	mux  *sync.RWMutex
}

type User struct {
	ID           int    `json:"id"`
	PasswordHash string `json:"password"`
	Email        string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}
