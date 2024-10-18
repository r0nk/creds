package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"os/exec"
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
			add_cred(cred{username: strings.TrimSpace(parts[0]), password: strings.TrimSpace(parts[1])})
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func write_creds_to_file(file_path string) error {
	file, err := os.OpenFile(file_path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, c := range creds {
		_, err := fmt.Fprintf(file, "%s:%s\n", c.username, c.password)
		if err != nil {
			return err
		}
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
	}
}

func dump_users(creds []cred) {
	for _, a := range creds {
		fmt.Printf("%s\n", a.username)
	}
}

func dump_passwords(creds []cred) {
	for _, a := range creds {
		fmt.Printf("%s\n", a.password)
	}
}

func all_creds() {
	title_creds()
	capslock_creds()
	lowercase_creds()
	permutate_creds()
	dual_creds()
}

func select_creds(query string) []cred {
	var ret []cred
	for _, c := range creds {
		if strings.Contains(strings.ToLower(c.username), strings.ToLower(query)) {
			ret = append(ret, c)
			if one {
				break
			}
		}
	}
	if len(ret) == 0 { // no users? try passwords
		for _, c := range creds {
			if strings.Contains(strings.ToLower(c.password), strings.ToLower(query)) {
				ret = append(ret, c)
				if one {
					break
				}
			}
		}
	}
	return ret
}

//https://stackoverflow.com/questions/55300117/how-do-i-find-all-files-that-have-a-certain-extension-in-go-regardless-of-depth#67629473
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func run_cred_checks(creds []cred) {

	file_path, err := find_file("creds.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	dir := filepath.Dir(file_path)
	matches, err := WalkMatch(dir, "credcheck")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, match := range matches {
		for _, c := range creds {
			blue := "\033[34m"
			cyan := "\033[36m"
			yellow := "\033[33m"
			green := "\033[32m"
			reset := "\033[0m"
			fmt.Printf("%s Checking %s%s%s %swith (%s%s%s:%s%s%s) %s\n", blue, green, match, reset, blue, cyan, c.username, reset, yellow, c.password, blue, reset)
			output, err := exec.Command(match, c.username, c.password).CombinedOutput()
			if err != nil {
				fmt.Println("Running error:", err)
			}
			fmt.Printf("%s\n", output)
		}
	}
}

func main() {
	const filename = "creds.txt"

	var all_mods_flag = flag.Bool("M", false, "Apply all possible mutations to creds.")
	var only_passwords = flag.Bool("p", false, "Only dump passwords")
	var only_users = flag.Bool("u", false, "Only dump users")
	var one_flag = flag.Bool("1", false, "Only dump one result")
	var add_flag = flag.Bool("a", false, "Read credentials from stdin and add to creds.txt")
	var check_flag = flag.Bool("c", false, "Run every credcheck file under creds.txt directory.")
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

	if *add_flag {
		err = read_file_into_creds("/dev/stdin")
		if err != nil {
			fmt.Println("Error reading from STDIN:", err)
			return
		}
		err = write_creds_to_file(file_path)
		if err != nil {
			fmt.Println("Error writing to creds:", err)
			return
		}
		return
	}

	if *all_mods_flag {
		all_creds()
	}
	selected := select_creds(flag.Arg(0))
	if *check_flag {
		run_cred_checks(selected)
		return
	}
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
