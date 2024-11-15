package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadEnv() {
	filePath := ".env"
	environment := os.Getenv("ENVIRONMENT")

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	lines := make([]string, 0, 100)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	for _, l := range lines {
		pair := strings.Split(l, "=")
		os.Setenv(pair[0], pair[1])
	}

	if environment == "SERVICE_DEV" {
		os.Setenv("POSTGRES_DB_HOST", os.Getenv("POSTGRES_DB_SERVICE_HOST"))
	}
}
