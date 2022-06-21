package main

import "os"

type FileDb struct {
	file *os.File
}

// func New(config Config) {
// 	file, err := os.OpenFile(config.Address)
// }
