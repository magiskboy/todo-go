package main

import "os"

func main() {
	os.Setenv("ACCESS_SECRET", "lmao")
	InitDB()
	InitHTTP()
	App.Run()
}
