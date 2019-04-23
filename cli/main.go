package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"shuSemester/infrastructure"
	"shuSemester/service/token"
)

func main() {
	username := os.Args[1]
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(username), -1)
	_, err := infrastructure.DB.Exec(`
	INSERT INTO token(tokenhash) VALUES ($1);
	`, string(encrypted))
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(token.GenerateJWT(username))
}
