package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestCdkStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkStack(app, "MyStack", nil)

	// THEN
	template := assertions.Template_FromStack(stack, nil)

	// Check for VPC resource
	template.HasResourceProperties(jsii.String("AWS::EC2::VPC"), map[string]interface{}{
		"CidrBlock":          jsii.String(vpcCidrDefault),
		"EnableDnsHostnames": jsii.Bool(true),
		"EnableDnsSupport":   jsii.Bool(true),
	})

	// Check for Subnet resource
	template.HasResourceProperties(jsii.String("AWS::EC2::Subnet"), map[string]interface{}{
		"CidrBlock":           jsii.String("10.0.1.0/24"),
		"MapPublicIpOnLaunch": jsii.Bool(true),
	})
	// Check for Subnet resource
	template.HasResourceProperties(jsii.String("AWS::EC2::Subnet"), map[string]interface{}{
		"CidrBlock":           jsii.String("10.0.0.0/24"),
		"MapPublicIpOnLaunch": jsii.Bool(true),
	})

	// Check for Internet Gateway
	template.HasResourceProperties(jsii.String("AWS::EC2::InternetGateway"), map[string]interface{}{})

	// Check for Route Table
	template.HasResourceProperties(jsii.String("AWS::EC2::RouteTable"), map[string]interface{}{})

	// Check for Route
	template.HasResourceProperties(jsii.String("AWS::EC2::Route"), map[string]interface{}{
		"DestinationCidrBlock": jsii.String("0.0.0.0/0"),
	})

	// Check for Security Group
	template.HasResourceProperties(jsii.String("AWS::EC2::SecurityGroup"), map[string]interface{}{
		"GroupDescription": jsii.String("Allow access to EKS Cluster"),
	})

	// Check for EKS Cluster resource
	template.HasResourceProperties(jsii.String("Custom::AWSCDK-EKS-Cluster"), map[string]interface{}{})

	// Check for Node Group
	template.HasResourceProperties(jsii.String("AWS::EKS::Nodegroup"), map[string]interface{}{
		"ScalingConfig": map[string]interface{}{
			"MinSize":     jsii.Number(minNodeSizeDefault),
			"MaxSize":     jsii.Number(maxNodeSizeDefault),
			"DesiredSize": jsii.Number(desiredNodeSizeDefault),
		},
		"DiskSize":      jsii.Number(diskSizeDefault),
		"InstanceTypes": []interface{}{nodeInstanceTypeDefault},
	})
}
