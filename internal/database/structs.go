package database

import (
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RefreshTokens map[int]RefreshToken `json:"refresh_tokens"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID           int    `json:"id"`
	PasswordHash string `json:"password"`
	Email        string `json:"email"`
}

type RefreshToken struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}
