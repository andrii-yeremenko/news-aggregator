package main

import (
	"flag"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"os"
)

const (
	userNameDefault          = "andrii"
	vpcCidrDefault           = "10.0.0.0/16"
	subnetsMaskDefault       = 24
	eksClusterNameDefault    = "news-aggregator-cluster"
	kubernetesVersionDefault = "1.30"
	nodeInstanceTypeDefault  = "t3.medium"
	minNodeSizeDefault       = 1
	maxNodeSizeDefault       = 3
	desiredNodeSizeDefault   = 2
	ec2SshKeyDefault         = "Main"
	diskSizeDefault          = 20
)

var addonVersions = map[string]string{
	"1.28:VpcCniAddonVersion":      "v1.18.3-eksbuild.3",
	"1.28:KubeProxyAddonVersion":   "v1.28.12-eksbuild.5",
	"1.28:CoreDnsAddonVersion":     "v1.10.1-eksbuild.13",
	"1.28:PodIdentityAddonVersion": "v1.2.0-eksbuild.1",
	"1.29:VpcCniAddonVersion":      "v1.18.3-eksbuild.3",
	"1.29:KubeProxyAddonVersion":   "v1.29.7-eksbuild.5",
	"1.29:CoreDnsAddonVersion":     "v1.11.3-eksbuild.1",
	"1.29:PodIdentityAddonVersion": "v1.3.2-eksbuild.2",
	"1.30:VpcCniAddonVersion":      "v1.18.3-eksbuild.3",
	"1.30:KubeProxyAddonVersion":   "v1.30.3-eksbuild.5",
	"1.30:CoreDnsAddonVersion":     "v1.11.3-eksbuild.1",
	"1.30:PodIdentityAddonVersion": "v1.3.2-eksbuild.2",
}

var (
	userName          string
	vpcCidr           string
	subnetsMask       float64
	eksClusterName    string
	kubernetesVersion string
	nodeInstanceType  string
	minNodeSize       float64
	maxNodeSize       float64
	desiredNodeSize   float64
	ec2SshKey         string
	diskSize          float64
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	fetchParams(stack)

	vpc := awsec2.NewVpc(stack, jsii.String(getPrefixedName("vpc")), &awsec2.VpcProps{
		IpAddresses: awsec2.IpAddresses_Cidr(jsii.String(vpcCidr)),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:                jsii.String(getPrefixedName("public-subnet")),
				SubnetType:          awsec2.SubnetType_PUBLIC,
				CidrMask:            jsii.Number(subnetsMask),
				MapPublicIpOnLaunch: jsii.Bool(true),
			},
		},
		MaxAzs: jsii.Number(2),
	})

	eksSecurityGroup := awsec2.NewSecurityGroup(stack, jsii.String(getPrefixedName("eks-sg")), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		Description:       jsii.String("Allow access to EKS Cluster"),
		SecurityGroupName: jsii.String(getPrefixedName("eks-sg")),
	})

	eksSecurityGroup.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.Port_AllTraffic(),
		jsii.String("Allow all inbound traffic"),
		jsii.Bool(false),
	)

	eksRole := awsiam.NewRole(stack, jsii.String(getPrefixedName("eks-role")), &awsiam.RoleProps{
		AssumedBy:   awsiam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), nil),
		Description: jsii.String("Role for EKS Cluster"),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSClusterPolicy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSServicePolicy")),
		},
		RoleName: jsii.String(getPrefixedName("eks-role")),
	})

	eksCluster := awseks.NewCluster(stack, jsii.String(getPrefixedName("eks-cluster")), &awseks.ClusterProps{
		Version:       awseks.KubernetesVersion_Of(jsii.String(kubernetesVersion)),
		ClusterName:   jsii.String(eksClusterName),
		Role:          eksRole,
		SecurityGroup: eksSecurityGroup,
		Vpc:           vpc,
		VpcSubnets: &[]*awsec2.SubnetSelection{
			{Subnets: &[]awsec2.ISubnet{(*vpc.PublicSubnets())[0]}},
			{Subnets: &[]awsec2.ISubnet{(*vpc.PublicSubnets())[1]}},
		},
		EndpointAccess:  awseks.EndpointAccess_PUBLIC(),
		IpFamily:        "IP_V4",
		DefaultCapacity: jsii.Number(0),
	})

	iamUserArn := "arn:aws:iam::" + os.Getenv("CDK_DEFAULT_ACCOUNT") + ":user/" + userName
	eksCluster.AwsAuth().AddUserMapping(awsiam.User_FromUserArn(stack, jsii.String(userName), jsii.String(iamUserArn)), &awseks.AwsAuthMapping{
		Username: jsii.String(userName),
		Groups: &[]*string{
			jsii.String("system:masters"),
		},
	})

	nodeGroupRole := awsiam.NewRole(stack, jsii.String(getPrefixedName("node-group-role")), &awsiam.RoleProps{
		AssumedBy:   awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
		Description: jsii.String("Role for EKS Node Group"),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ContainerRegistryReadOnly")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2FullAccess")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKS_CNI_Policy")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSWorkerNodePolicy")),
		},
		RoleName: jsii.String(getPrefixedName("node-group-role")),
	})

	eksCluster.AddNodegroupCapacity(jsii.String(getPrefixedName("node-group")), &awseks.NodegroupOptions{
		InstanceTypes: &[]awsec2.InstanceType{
			awsec2.NewInstanceType(jsii.String(nodeInstanceType)),
		},
		NodeRole:    nodeGroupRole,
		MinSize:     jsii.Number(minNodeSize),
		MaxSize:     jsii.Number(maxNodeSize),
		DesiredSize: jsii.Number(desiredNodeSize),
		RemoteAccess: &awseks.NodegroupRemoteAccess{
			SshKeyName: jsii.String(ec2SshKey),
		},
		Subnets: &awsec2.SubnetSelection{
			Subnets: vpc.PublicSubnets(),
		},
		AmiType:  awseks.NodegroupAmiType_AL2_X86_64,
		DiskSize: jsii.Number(diskSize),
	})

	addOutput(stack, "EKS Cluster Name", "EksClusterName", eksCluster.ClusterName())
	addOutput(stack, "EKS Cluster Endpoint", "EksClusterEndpoint", eksCluster.ClusterEndpoint())
	addOutput(stack, "EKS Cluster Security Group", "EksClusterSecurityGroup", eksCluster.ClusterSecurityGroup().SecurityGroupId())
	addOutput(stack, "EKS Cluster Role", "EksClusterRole", eksCluster.Role().RoleArn())
	addOutput(stack, "EKS Node Group Role", "EksNodeGroupRole", nodeGroupRole.RoleArn())
	addOutput(stack, "VPC add-on version", "VpcCniAddonVersion", jsii.String(getAddonVersion(kubernetesVersion, "VpcCniAddonVersion")))
	addOutput(stack, "Kube Proxy add-on version", "KubeProxyAddonVersion", jsii.String(getAddonVersion(kubernetesVersion, "KubeProxyAddonVersion")))
	addOutput(stack, "CoreDNS add-on version", "CoreDnsAddonVersion", jsii.String(getAddonVersion(kubernetesVersion, "CoreDnsAddonVersion")))
	addOutput(stack, "Pod Identity add-on version", "PodIdentityAddonVersion", jsii.String(getAddonVersion(kubernetesVersion, "PodIdentityAddonVersion")))

	return stack
}

func fetchParams(stack awscdk.Stack) {
	getString := func(key string, defaultValue string) string {
		if value := stack.Node().TryGetContext(jsii.String(key)); value != nil {
			if str, ok := value.(string); ok {
				return str
			}
		}
		return defaultValue
	}

	getFloat64 := func(key string, defaultValue float64) float64 {
		if value := stack.Node().TryGetContext(jsii.String(key)); value != nil {
			if num, ok := value.(float64); ok {
				return num
			}
		}
		return defaultValue
	}

	userName = getString("userName", userNameDefault)
	vpcCidr = getString("vpcCidr", vpcCidrDefault)
	subnetsMask = getFloat64("subnetsMask", float64(subnetsMaskDefault))
	eksClusterName = getString("eksClusterName", eksClusterNameDefault)
	kubernetesVersion = getString("kubernetesVersion", kubernetesVersionDefault)
	nodeInstanceType = getString("nodeInstanceType", nodeInstanceTypeDefault)
	minNodeSize = getFloat64("minNodeSize", float64(minNodeSizeDefault))
	maxNodeSize = getFloat64("maxNodeSize", float64(maxNodeSizeDefault))
	desiredNodeSize = getFloat64("desiredNodeSize", float64(desiredNodeSizeDefault))
	ec2SshKey = getString("ec2SshKey", ec2SshKeyDefault)
	diskSize = getFloat64("diskSize", float64(diskSizeDefault))
}

func addOutput(stack awscdk.Stack, desc string, displayName string, value *string) {
	awscdk.NewCfnOutput(stack, jsii.String(desc), &awscdk.CfnOutputProps{
		Description: jsii.String(desc),
		Value:       value,
		ExportName:  jsii.String(displayName),
	})
}

func getAddonVersion(k8sVersion string, addonName string) string {
	key := k8sVersion + ":" + addonName
	if version, ok := addonVersions[key]; ok {
		return version
	}

	availableK8sVersions := make([]string, 0, len(addonVersions))
	for k := range addonVersions {
		availableK8sVersions = append(availableK8sVersions, k)
	}

	err := "Addon version compatible for this Kubernetes version not found! Available versions: "
	for _, v := range availableK8sVersions {
		err += v + " "
	}

	panic(err)
}

func getPrefixedName(name string) string {
	return userName + "-" + name
}

func main() {
	defer jsii.Close()

	flag.Parse()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "andrii-cdk-stack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
