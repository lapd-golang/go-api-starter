package models

func CheckAuth(username string, password string) bool {
	user := User{}
	user = user.GetByUsernameAndPassword(username, password)

	if user.ID > 0 {
		return true
	}

	return false
}