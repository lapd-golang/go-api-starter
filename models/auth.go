package models

func CheckAuth(username string, password string) bool {
	user := User{
		Username:username,
		Password:password,
	}
	user = user.CheckUser()

	if user.ID > 0 {
		return true
	}

	return false
}