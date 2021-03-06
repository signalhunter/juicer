package main

import (
	"bufio"
	"os"
)

func ReadFile(path string) ([]string, error) {
	data := make([]string, 0)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) > 0 && txt[0] != '#' {
			data = append(data, txt)
		}
	}
	return data, nil
}

func WriteResults(path string, recv chan string, done chan bool) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for d := range recv {
		writer.WriteString(d + "\n")
	}
	writer.Flush()
	done <- true
	return nil
}
