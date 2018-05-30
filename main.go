package main

func main() {
	app := App{}
	app.Initialize(
		"192.168.99.100",
		"godb",
		"movies",
	)

	app.Run("27017")
}

