package utility

func CheckError(e error, msg string) {
	if e != nil {
		panic(msg + e.Error())
	}
}