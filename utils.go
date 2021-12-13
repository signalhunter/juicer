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
		data = append(data, scanner.Text())
	}
	return data, nil
}

func WriteResults(path string, recv chan string, done chan bool) error {
	dedup := make(map[string]bool)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for d := range recv {
		id := d[len(d)-11:]
		if !dedup[id] {
			dedup[id] = true
			writer.WriteString(id + "\n")
		}
	}
	writer.Flush()
	done <- true
	return nil
}
