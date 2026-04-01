package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	content, err := os.ReadFile("apps/server/go.mod")
	if err != nil {
		panic(err)
	}

	str := string(content)
	str = strings.Replace(str, "go 1.25.0", "go 1.24.3", 1) // Force 1.24.3 in case it was auto bumped

	err = os.WriteFile("apps/server/go.mod", []byte(str), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done")
}
