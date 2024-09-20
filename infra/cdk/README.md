# News Aggregator AWS CDK

This project uses the AWS Cloud Development Kit (CDK) to define and deploy infrastructure for a news aggregator
application. It sets up an Amazon EKS (Elastic Kubernetes Service) cluster with necessary networking and IAM
configurations.

## Prerequisites

- [AWS CLI](https://aws.amazon.com/cli/) configured with your credentials
- [Node.js](https://nodejs.org/en/) (for AWS CDK)
- [Go](https://golang.org/dl/) (for running the application)

## Getting Started

1. Clone the repository and navigate to the `infra/cdk` directory.
2. Install the required dependencies by running `npm install`.
```shell
$ npm install
```

3. Bootstrap the CDK toolkit by running `cdk bootstrap`.
```shell
$ cdk bootstrap
```

4. Deploy the stack by running `cdk deploy`.
```shell
$ cdk deploy
```

5. Once the stack is deployed, you can access the Kubernetes cluster by running `kubectl` commands.
```shell
$ aws eks --region $(aws configure get region) update-kubeconfig --name $(cdk context cluster-name)
$ kubectl get nodes
```

6. To delete the stack, run `cdk destroy`.
```shell
$ cdk destroy
```

## Configuration

CloudFormation parameters can be configured in the `cdk.json` file. The following parameters are available:

```json
{
  "app": "go run main.go",
  "context": {
    "userName": "your-username",
    "vpcCidr": "10.0.0.0/16",
    "subnetsMask": 24,
    "eksClusterName": "your-cluster-name",
    "kubernetesVersion": "1.30",
    "nodeInstanceType": "t3.medium",
    "minNodeSize": 1,
    "maxNodeSize": 3,
    "desiredNodeSize": 2,
    "ec2SshKey": "your-ssh-key",
    "diskSize": 20
  }
}
```
This CDK stack supports `1.28`, `1.29`, and `1.30` Kubernetes versions. The default version is `1.30`.

## Outputs

After deployment, the following outputs will be available:

* EKS Cluster Name
* EKS Cluster Endpoint
* EKS Cluster Security Group
* EKS Cluster Role
* EKS Node Group Role
* VPC add-on versions (for CoreDNS, Kube Proxy, etc.)

## License

This project is licensed under the TeamDev License