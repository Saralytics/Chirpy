package database

import "sync"

type DB struct {
	path   string
	mux    *sync.RWMutex
	chirps map[int]Chirp
	nextID int
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}
