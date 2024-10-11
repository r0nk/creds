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

var one bool

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

func dump_creds(creds []cred) {
	for _, a := range creds {
		fmt.Printf("%s:%s\n", a.username, a.password)
		if one {
			return
		}
	}
}

func dump_users(creds []cred) {
	for _, a := range creds {
		fmt.Printf("%s\n", a.username)
		if one {
			return
		}
	}
}

func dump_passwords(creds []cred) {
	for _, a := range creds {
		fmt.Printf("%s\n", a.password)
		if one {
			return
		}
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

func select_creds(query string) []cred {
	var ret []cred
	for _, c := range creds {
		if strings.Contains(strings.ToLower(c.username), strings.ToLower(query)) {
			ret = append(ret, c)
		}
	}
	if len(ret) == 0 { // no users? try passwords
		for _, c := range creds {
			if strings.Contains(strings.ToLower(c.password), strings.ToLower(query)) {
				ret = append(ret, c)

			}
		}
	}
	return ret
}

func main() {
	const filename = "creds.txt"

	var all_mods_flag = flag.Bool("M", false, "Apply all possible mutations to creds.")
	var only_passwords = flag.Bool("p", false, "Only dump passwords")
	var only_users = flag.Bool("u", false, "Only dump users")
	var one_flag = flag.Bool("1", false, "Only dump one result")
	flag.Parse()

	one = *one_flag

	if *only_passwords && *only_users {
		fmt.Printf("Error: -u and -p are exlusive options")
		return
	}

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
	selected := select_creds(flag.Arg(0))
	if !*only_users && !*only_passwords {
		dump_creds(selected)
	}
	if *only_users {
		dump_users(selected)
	}

	if *only_passwords {
		dump_passwords(selected)
	}
}
