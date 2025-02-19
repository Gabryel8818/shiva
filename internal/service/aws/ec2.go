package ec2

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Instance struct {
	ID           string
	InstanceType string
	Name         string
}

var cfg aws.Config

func init() {
	var err error
	cfg, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Error loading configuration: %v", err))
	}
}

func GetInstances() ([]Instance, error) {
	ec2Client := ec2.NewFromConfig(cfg)

	var taggedEc2Instances []Instance

	var nextToken *string

	for {

		instancesOutput, err := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{
			MaxResults: aws.Int32(100),
			NextToken:  nextToken,
		})

		if err != nil {
			return nil, fmt.Errorf("error getting ec2 instances: %v", err)
		}

		for _, reservation := range instancesOutput.Reservations {
			for _, instance := range reservation.Instances {
				taggedEc2Instance := Instance{
					ID:           *instance.InstanceId,
					InstanceType: string(instance.InstanceType),
				}

				shivaManaged := hasShivaManagedTag(instance.Tags)

				if shivaManaged {
					taggedEc2Instance.Name = getInstanceTagValue(instance.Tags, "Name")

					taggedEc2Instances = append(taggedEc2Instances, taggedEc2Instance)
				}

			}
		}

		nextToken = instancesOutput.NextToken

		if nextToken == nil {
			break
		}
	}

	return taggedEc2Instances, nil
}

func hasShivaManagedTag(tags []types.Tag) bool {
	for _, tag := range tags {
		if *tag.Key == "shiva:managed" {
			return true
		}

	}
	return false
}

func getInstanceTagValue(tags []types.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

func StopEC2Instances(instances []Instance) (string, error) {
	var wg sync.WaitGroup
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", fmt.Errorf("falha ao carregar a configuração da AWS: %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	for _, instance := range instances {
		wg.Add(1)
		go func(instance Instance) {
			defer wg.Done()

			input := &ec2.StopInstancesInput{
				InstanceIds: []string{instance.ID},
			}

			_, err := client.StopInstances(context.Background(), input)
			if err != nil {
				log.Printf("Erro ao desligar a instância %s: %v", instance.ID, err)
			} else {
				log.Printf("Instância %s desligada com sucesso!", instance.ID)
			}
		}(instance)
	}

	// Aguarda todas as goroutines terminarem
	wg.Wait()

	return "Instances stopped", nil
}
