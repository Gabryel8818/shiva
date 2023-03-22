package main

import (
	"fmt"
	"log"
	"shiva/ec2"
)

func main() {
	fmt.Println("Hi, I'm Shiva.")

	instances, err := ec2.GetInstances()
	if err != nil {
		log.Fatalf("Error getting instances: %v", err)
	}

	fmt.Println(instances)
}
