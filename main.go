package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func main() {
	fmt.Println("Hi, im shiva.")

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Error loading configuration")
		return
	}

	ec2Client := ec2.NewFromConfig(cfg)

	instances, err := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
	if err != nil {
		fmt.Println("Error getting ec2 instances", err)
		return
	}

	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println("Instance ID", *instance.InstanceId)
			fmt.Println("Instance Type", string(instance.InstanceType))
			fmt.Println("Instance State", string(instance.State.Name))
		}
	}
}
