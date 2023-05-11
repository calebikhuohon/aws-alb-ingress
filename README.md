# Setup

This demo uses Kubernetes as its compute option, and Terraform as the infrastructure as code tool of choice.

The ALB Ingress controller is used to handle network termination for the Kubernetes cluster

## Prerequisites
* [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) installed and configured for the AWS account in use.
* [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)

## Running The Application
The application can be deployed via bash scripts or GitHub Actions which handles the infrastructure provisioning and application deployment.

The scripts include:
* `provision-infra.sh` script which creates the Kubernetes cluster, its ingress controller setup and ecr repository.
* `deploy-app.sh` which deploys the application to kubernetes and handles ingress creation.

### Steps

#### GitHub Actions
* Set your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` in GitHub Secrets
* Run the various workflows in the GitHub Actions tab

#### Command Line Scripts
* `provision-infra.sh`
  * Replace the `AWS_USER_ARN` and `AWS_ARN` which your AWS account's equivalent.
  * Optional: change the `AWS_REGION` and `CLUSTER_NAME` to preferred values.
  * Run `chmod +x provision-infra.sh` to make the script executable.
  * Run `./provision-infra.sh` to execute the script. It will require a few minutes for completion.
  

* `deploy-app.sh`
  * Change the `ECR_REGISTRY` variable to the value created for your AWS account
  * Optional: change `AWS_REGION` and `CLUSTER_NAME` to preferred values
  * Run `chmod +x deploy-app.sh` to make the script executable.
  * Run `./deploy-app.sh` to execute the script. 

* Open the AWS management console to get the auto-generated [load balancers](https://us-east-2.console.aws.amazon.com/ec2/home?region=us-east-2#LoadBalancers) DNS name

* Enter `http://<LOAD BALANCER URL>/get-res` in a browser to access the application's endpoint and retrieve the contents of the file uploaded to S3 
  e.g http://k8s-default-restserv-a18b9b9a6b-419094574.us-east-2.elb.amazonaws.com/get-res

  This returns `{
  "greeting": "I am the Foo"
  }`

