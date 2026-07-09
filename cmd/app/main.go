package main

import ()

func main() {
	application := app{
		port: ":8080",
	}

	application.serve(application.mount())
}
