package token

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"shuSemester/infrastructure"
)

func ValidateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return false
	}
	claims := token.Claims.(jwt.MapClaims)
	studentId := claims["studentId"].(string)
	rows, _ := infrastructure.DB.Query(`
	SELECT tokenHash from token;
	`)
	for rows.Next() {
		var thisToken string
		_ = rows.Scan(&thisToken)
		if bcrypt.CompareHashAndPassword([]byte(thisToken), []byte(studentId)) == nil {
			return true
		}
	}
	return false
}

func GenerateJWT(id string) string {
	result, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"studentId": id,
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	return result
}
