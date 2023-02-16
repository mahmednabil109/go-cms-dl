package utils

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// Dotfile handles parsing simple key value dot files and exports
// a function to query these values
type Dotfile struct {
	table map[string]string
}

// NewDotfile constructes a new dotfile struct from file in the path
// file should be in the form of key=value
func NewDotfile(path string) (*Dotfile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dotfile := Dotfile{
		table: make(map[string]string),
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \n")
		row := strings.Split(line, "=")
		if len(row) != 2 {
			return nil, errors.New("failed to parse line " + line)
		}
		dotfile.table[row[0]] = row[1]
	}
	return &dotfile, nil
}

// Get returns the value associated with `key` in the dotfile
func (f *Dotfile) Get(key string) string {
	if f.table == nil {
		return ""
	}
	return f.table[key]
}
