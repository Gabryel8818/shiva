package ec2

import (
	"context"
	"fmt"

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

type Tag struct {
	Key   string
	Value string
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

	instancesOutput, err := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("error getting ec2 instances: %v", err)
	}

	var taggedEc2Instances []Instance

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
