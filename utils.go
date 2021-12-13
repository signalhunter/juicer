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

func WriteResults(path string, recv chan string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for d := range recv {
		writer.WriteString(d + "\n")
	}
	defer writer.Flush()
	return nil
}