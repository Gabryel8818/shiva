package main

import (
	"fmt"
	"log"
	ec2 "shiva/internal/service/aws"
)

func main() {
	fmt.Println("Hi, I'm Shiva.")

	instances, err := ec2.GetInstances()
	if err != nil {
		log.Fatalf("Error getting instances: %v", err)
	}

	fmt.Println(instances)

	result, err := ec2.StopEC2Instances(instances)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Println(result)

}
