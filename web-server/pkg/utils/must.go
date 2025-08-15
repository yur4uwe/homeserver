package utils

func Must(value interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return value
}
