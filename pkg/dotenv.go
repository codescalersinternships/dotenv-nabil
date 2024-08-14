package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Load will read your env file(s) and load them into ENV for this process.
//
// Call this function as close as possible to the start of your program (ideally in main).
//
// If you call Load without any args it will default to loading .env in the current path.
//
// You can otherwise tell it which files to load (there can be more than one) like:
//
//	godotenv.Load("fileone", "filetwo")
//
// It's important to note that it WILL NOT OVERRIDE an env variable that already exists - consider the .env file to set dev vars or sensible defaults.
func Load(filenames ...string) (err error) {
	if len(filenames) == 0 {
		return loadFile(".env")
	}
	for _, filename := range filenames {
		err = loadFile(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

// Parser reads an env file from io.Reader, returning a map of keys and values.
func Parser(reader io.Reader) (map[string]string, error) {
	out := make(map[string]string)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if !strings.ContainsAny(line, "=") && !strings.ContainsAny(line, ":") {
			return nil, fmt.Errorf("line is not a key value pair")
		}
		var keyVal []string
		if strings.ContainsAny(line, "=") {
			keyVal = strings.SplitN(line, "=", 2)
		} else {
			keyVal = strings.SplitN(line, ":", 2)
		}
		var key, val string
		key = keyVal[0]
		val = keyVal[1]
		if len(line) >= 6 && line[:6] == "export" {
			key = key[6:]
		}
		if strings.ContainsAny(val, "#") {
			val = val[:strings.Index(val, "#")]
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if len(key) == 0 || len(val) == 0 {
			return nil, fmt.Errorf("key val aren't valid")
		}
		out[key] = val
	}
	return out, nil
}

func loadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	in := bufio.NewReader(file)
	out, err := Parser(in)
	if err != nil {
		return err
	}

	for key, value := range out {
		err = os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// Unmarshal reads an env file from a string, returning a map of keys and values.
func Unmarshal(str string) (map[string]string, error) {
	return Parser(strings.NewReader(str))
}
