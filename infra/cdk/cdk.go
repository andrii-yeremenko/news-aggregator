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

// To run this CDK application, first define environment variables in cdk.json file.
// Otherwise, the application will use default values for the parameters.
// The following default values for the parameters are:
const (
	// Default username that will be used for the resource prefix
	userNameDefault = "andrii"
	// Default CIDR block for the VPC.
	// This is the IP address range that will be used for the VPC.
	// As a default, it is class A network.
	vpcCidrDefault = "10.0.0.0/16"
	// Default mask for the subnets.
	// This will affect the number of available IP addresses in the subnet.
	// Recommend you to use it at least /26 for the subnets.
	// The default value is /24.
	subnetsMaskDefault = 24
	// Default name for the EKS cluster. This is the name of the EKS cluster that will be created.
	eksClusterNameDefault = "news-aggregator-cluster"
	// Default Kubernetes version. Check CDK README for the list of available versions.
	kubernetesVersionDefault = "1.30"
	// Default instance type for the nodes. This is the type of EC2 instances that will be created for the nodes.
	nodeInstanceTypeDefault = "t3.medium"
	// Default minimum number of nodes. This is the minimum number of nodes that can be created.
	minNodeSizeDefault = 1
	// Default maximum number of nodes. This is the maximum number of nodes that can be created.
	maxNodeSizeDefault = 3
	// Default desired number of nodes. This is the number of nodes that will be created at the beginning.
	desiredNodeSizeDefault = 2
	// Default EC2 SSH key that will be used for the instances to access them in the future.
	ec2SshKeyDefault = "Main"
	// Default disk size for the nodes
	diskSizeDefault = 20
)

// This is a map of add-on versions for different Kubernetes versions.
// The key is a combination of Kubernetes version and add-on name, and the value is the version of the add-on.
// Update this map if you want to use a different version of the add-on for a specific Kubernetes version.
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

// Parameters for the stack
// See the CDK context in the cdk.json file for the default values
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

// CdkStackProps defines the properties for the stack
type CdkStackProps struct {
	awscdk.StackProps
}

// NewCdkStack creates a new CDK stack
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

// fetchParams fetches the parameters from the CDK context
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

// addOutput adds a single output to the stack
func addOutput(stack awscdk.Stack, desc string, displayName string, value *string) {
	awscdk.NewCfnOutput(stack, jsii.String(desc), &awscdk.CfnOutputProps{
		Description: jsii.String(desc),
		Value:       value,
		ExportName:  jsii.String(displayName),
	})
}

// getAddonVersion returns the version of the add-on compatible with the Kubernetes version
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

// getPrefixedName returns a name with a prefix, for example "<prefix>-<resource-name>"
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

// env returns the environment for the stack
func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
