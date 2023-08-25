package main

import (
	"fmt"

	"github.com/lcox74/tundra-dns/backend/internal/routing"
)

func main() {
	fmt.Println("Hello, World!")

	// Launch the DNS Query Handler
	routing.LaunchDNSQueryHandler()
}
