package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "Migrations are not configured: select and configure a database first.")
	os.Exit(1)
}
