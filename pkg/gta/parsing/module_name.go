package parsing

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path"
)

func GetModuleName() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	goModPath := path.Join(dir, "go.mod")

	file, err := os.Open(goModPath)
	if err != nil {
		return ""
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := getFullLine(reader)
		if err != nil && err != io.EOF {
			return ""
		}

		if bytes.Contains(line, []byte("module")) {
			return getModuleName(line)
		}

		if err != nil && err == io.EOF {
			return "nil"
		}
	}
}

func getModuleName(line []byte) string {
	lineSplit := bytes.SplitN(line, []byte("module"), 2)
	if len(lineSplit) != 2 {
		return ""
	}

	moduleName := string(bytes.TrimSpace(bytes.TrimSpace(lineSplit[1])))

	return moduleName
}

func getFullLine(reader *bufio.Reader) ([]byte, error) {
	var (
		line     []byte
		tmpLine  []byte
		isPrefix bool
		err      error
	)

	for {
		tmpLine, isPrefix, err = reader.ReadLine()
		if err == io.EOF {
			return line, io.EOF
		} else if err != nil {
			return nil, err
		}

		line = append(line, tmpLine...)
		if !isPrefix {
			break
		}
	}

	return line, nil
}
