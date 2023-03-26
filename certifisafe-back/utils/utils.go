package utils

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckSQLError(err error) {
	if err != nil {
		panic("SQL error")
	}
}
