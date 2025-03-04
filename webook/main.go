package main

func main() {
	server := InitWebServer()
	server.Run("0.0.0.0:8081")
}
