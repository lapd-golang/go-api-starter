package models

import "fmt"

func CheckAuth(username string, password string) bool {
	user := User{}
	user = user.GetByUsernameAndPassword(username, password)

	fmt.Printf("%+v\n", user)
	fmt.Printf("%s\n", user.CreatedAt)

	if user.ID > 0 {
		return true
	}

	return false
}