package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func readDB() []URL {
	jsonFile, err := os.Open("./db.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var urls []URL

	err = json.Unmarshal(byteValue, &urls)
	if err != nil {
		fmt.Println("Error unmarshalling JSON", err)
		return nil
	}

	return urls
}

func generateRandomCode() string {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
