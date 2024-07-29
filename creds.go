package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
)

type cred struct {
	username string
	password string
}

func find_file(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		file_path := filepath.Join(dir, filename)
		if _, err := os.Stat(file_path); err == nil {
			return file_path, nil
		}

		parent_dir := filepath.Dir(dir)
		if parent_dir == dir {
			break
		}
		dir = parent_dir
	}

	return "", fmt.Errorf("Credential file %s not found.", filename)
}

func read_file_into_creds(file_path string) ([]cred, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var creds []cred
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			creds = append(creds, cred{username: strings.TrimSpace(parts[0]), password: strings.TrimSpace(parts[1])})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return creds, nil
}

func title_creds(creds []cred) {
	for _, c := range creds {
		fmt.Printf("%s:%s\n", cases.Title(language.Und, cases.NoLower).String(c.username), c.password)
	}
}

func main() {
	const filename = "creds.txt"

	file_path, err := find_file(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	creds, err := read_file_into_creds(file_path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	command := os.Args[1]
	switch command {
	case "usernames":
		if len(os.Args) < 3 {
			for _, p := range creds {
				fmt.Printf("%s\n", p.username)
			}
		} else {
			password := os.Args[2]
			for _, p := range creds {
				if p.password == password {
					fmt.Printf("%s\n", p.username)
				}
			}
		}
	case "passwords":
		if len(os.Args) < 3 {
			for _, p := range creds {
				fmt.Printf("%s\n", p.username)
			}
		} else {
			username := os.Args[2]
			for _, p := range creds {
				if p.username == username {
					fmt.Printf("%s\n", p.password)
				}
			}
		}
	case "title":
		title_creds(creds)
	}
}
