package main

import (
	"crocus-calculator/internal/setting"
	"fmt"
	"os"
)

func main() {
	fmt.Println("=== Starting migrations ===")

	app := setting.App{}

	fmt.Println("1. Loading configuration...")
	app.LoadConfig()

	fmt.Println("2. Attempting to apply migrations...")
	err := app.Migrate()
	if err != nil {
		fmt.Printf("!!! Migration failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Migrations applied successfully")
	fmt.Println("=== Operation completed ===")
}
