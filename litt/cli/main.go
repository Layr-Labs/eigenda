package main

// main is the entry point for the LittDB cli.
func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
