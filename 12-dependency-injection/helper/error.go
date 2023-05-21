package helper

func PanifIfError(err error) {
	if err != nil {
		panic(err)
	}
}
