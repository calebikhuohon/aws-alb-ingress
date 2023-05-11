AWS_REGION=us-east-2
AWS_USER_ARN=arn:aws:iam::xxxx:user/xxxxx #replace with your AWS account user ARN
AWS_ARN=arn:aws:iam::xxxxxx #replace with your AWS account ARN
CLUSTER_NAME=ingress-setup-demo


#install eksctl
if ! command -v eksctl &> /dev/null
then
  curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
  sudo mv /tmp/eksctl /usr/local/bin
  echo eksctl version
fi

#install kubectl
if ! command -v kubectl  &> /dev/null
then
  curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
  chmod +x kubectl
  mkdir -p ~/.local/bin
  mv ./kubectl ~/.local/bin/kubectl
  echo kubectl version --client
fi

# install helm
if ! command -v helm &> /dev/null
then
  curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
  chmod 700 get_helm.sh
  ./get_helm.sh
fi

# create kubernetes cluster
eksctl create cluster --config-file=infrastructure/kubernetes/eksctl.yml

# setup kubeconfig
aws eks --region $AWS_REGION update-kubeconfig --name $CLUSTER_NAME


# setup cluster authentication
eksctl create iamidentitymapping \
    --cluster $CLUSTER_NAME \
    --region=$AWS_REGION \
    --arn "$AWS_USER_ARN" \
    --group eks-console-dashboard-restricted-access-group \
    --no-duplicate-arns

# setup ALB ingress controller
eksctl utils associate-iam-oidc-provider --region $AWS_REGION --cluster $CLUSTER_NAME --approve
curl -o iam-policy.json https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.4.7/docs/install/iam_policy.json
aws iam create-policy --policy-name AWSLoadBalancerControllerIAMPolicy --policy-document file://iam-policy.json
eksctl create iamserviceaccount \
--cluster=$CLUSTER_NAME \
--namespace=kube-system \
--name=aws-load-balancer-controller \
--attach-policy-arn=$AWS_ARN:policy/AWSLoadBalancerControllerIAMPolicy \
--override-existing-serviceaccounts \
--region $AWS_REGION \
--approve

kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.5.4/cert-manager.yaml
kubectl apply -f deployments/alb/aws-load-balancer.yml

####################################
# Terraform Resource Provisioning

cd infrastructure/terraform || exit

terraform fmt -check

terraform init

terraform validate -no-color

terraform plan -no-color -input=false

terraform apply -auto-approve -input=false

