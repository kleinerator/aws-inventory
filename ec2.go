package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
  "time"
  "encoding/json"
)

type SecurityGroup struct {
  GroupId               string `json:"group_id"`
  Name                  string `json:"name"`
}

type Tag struct {
  Key                   string `json:"key"`
  Value                 string `json:"value"`
}

type Instances struct {
  Instances             []Instance `json:"instances"`
}

type Instance struct {
  InstanceId            string `json:"instance_id"`
  ImageId               string `json:"image_id"`
  InstanceType          string `json:"instance_type"`
  LaunchTime            time.Time `json:"launch_time"`
  AvailabilityZone      string `json:"availability_zone"`
	Status								string `json:"status"`
  PublicDnsName         string `json:"public_dns_name"`
  PublicIpAddress       string `json:"public_ip_address"`
  PrivateDnsName        string `json:"private_dns_name"`
  PrivateIpAddress      string `json:"private_ip_address"`
  SecurityGroups        []SecurityGroup `json:"security_groups"`
  Tags                  []Tag `json:"tags"`
}

func main()  {
  region := "us-west-2"
  DescribeAllInstancesInRegion(region)
}

func DescribeAllInstancesInRegion(region string) {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})


	resp, err := svc.DescribeInstances(nil)
	//resp, err := svc.DescribeInstances(nil)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

  //fmt.Println(resp)
	instances := []Instance{}
  // resp has all of the response data, pull out instance IDs:
  // fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
  for idx, _ := range resp.Reservations {
      // fmt.Println("  > Number of instances: ", len(res.Instances))
      for _, inst := range resp.Reservations[idx].Instances {
         //fmt.Println("    - Instance ID: ", *inst.InstanceId)
         //fmt.Println(*inst)
         instance := Instance{
					 InstanceId: *inst.InstanceId,
					 ImageId: *inst.ImageId,
					 InstanceType: *inst.InstanceType,
					 LaunchTime: *inst.LaunchTime,
					 AvailabilityZone: *inst.Placement.AvailabilityZone,
					 Status: *inst.State.Name,
					 PrivateDnsName: *inst.PrivateDnsName,
					 PrivateIpAddress: *inst.PrivateIpAddress,
				 }

				 if *inst.State.Name != "stopped" {
	         instance.PublicDnsName = *inst.PublicDnsName
					 //fmt.Println("    - Instance PublicDnsName: ", *inst.PublicDnsName)
	         instance.PublicIpAddress = *inst.PublicIpAddress
					 //fmt.Println("    - Instance PublicIpAddress: ", *inst.PublicIpAddress)
			 	 }

         if len(inst.SecurityGroups) > 0 {
           for _, grp := range inst.SecurityGroups {
             group := SecurityGroup{}
             group.GroupId = *grp.GroupId
             group.Name = *grp.GroupName
             instance.SecurityGroups = append(instance.SecurityGroups, group)
           }
         }
         if len(inst.Tags) > 0 {
           for _, tg := range inst.Tags {
             tag := Tag{}
             tag.Key = *tg.Key
             tag.Value = *tg.Value
             instance.Tags = append(instance.Tags, tag)
           }
         }


         instances = append(instances, instance)
      }
  }

	fmt.Println(len(instances))
  contents, _ := json.MarshalIndent(instances, "", "  ")
  if err != nil {
    fmt.Printf("err")
  } else {
    fmt.Println(string(contents))
  }
}
