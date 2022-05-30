package main

import (
	"fmt"
	"log"
)

func main() {
	useCase := UseCase{}
	m, err := useCase.Influencer("@alexys_lozada")
	if err != nil {
		log.Printf("error: %v", err)
	}

	fmt.Println(m)
}
