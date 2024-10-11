package main

import (
	"bufio"
	"flag"
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

var creds []cred

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

func read_file_into_creds(file_path string) error {
	file, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			creds = append(creds, cred{username: strings.TrimSpace(parts[0]), password: strings.TrimSpace(parts[1])})
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func unique(cred cred, creds []cred) bool {
	for _, c := range creds {
		if c.username == cred.username && c.password == cred.password {
			return false
		}
	}
	return true
}

func add_cred(cred cred) {
	if unique(cred, creds) {
		creds = append(creds, cred)
	}
}

func title_creds() {
	for _, c := range creds {
		add_cred(cred{cases.Title(language.Und, cases.NoLower).String(c.username), c.password})
	}
}

func capslock_creds() {
	for _, c := range creds {
		add_cred(cred{strings.ToUpper(c.username), c.password})
	}
}

func lowercase_creds() {
	for _, c := range creds {
		add_cred(cred{strings.ToLower(c.username), c.password})
	}
}

func permutate_creds() {
	for _, a := range creds {
		for _, b := range creds {
			add_cred(cred{a.username, b.password})
		}
	}
}

func dual_creds() {
	for _, a := range creds {
		add_cred(cred{a.username, a.username})
	}
}

func dump_creds() {
	for _, a := range creds {
		fmt.Printf("%s:%s\n", a.username, a.password)
	}
}

func all_creds() {
	title_creds()
	capslock_creds()
	lowercase_creds()
	permutate_creds()
	dual_creds()
}

func check_creds(creds []cred) {
	//Recursively get services below the creds.txt file
	//directory structure like username/service/results.txt // smith/ftp/attempt.txt
	//for c := range checks{
	//	results := c(cred)
	//	TODO write results to correct filename
	//}
}

func main() {
	const filename = "creds.txt"

	var all_mods_flag = flag.Bool("M", false, "Apply all possible mutations to creds.")
	flag.Parse()

	file_path, err := find_file(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = read_file_into_creds(file_path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if *all_mods_flag {
		all_creds()
	}
	dump_creds()
}
