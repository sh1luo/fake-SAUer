package main

func init() {
	err := ReadConfig()
	if err != nil {
		panic(err)
	}
}
