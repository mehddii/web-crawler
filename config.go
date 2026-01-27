package main

import (
	"fmt"
	"os"
)

func GetConnStr() string {
	return fmt.Sprintf(
		`postgres://%s:%s@database/%s?sslmode=%s`,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("SSLMODE"),
	)
}
