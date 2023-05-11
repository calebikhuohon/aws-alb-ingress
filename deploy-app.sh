ECR_REGISTRY=450009495847.dkr.ecr.us-east-2.amazonaws.com #change to the ecr registry created for your AWS account
ECR_REPOSITORY=backend
IMAGE_TAG=latest
AWS_REGION=us-east-2
CLUSTER_NAME=DevopsTest
BUCKET_NAME=s3://ikh-json-bucket #If a new bucket was created, please replace with the bucket name
FILE_NAME=file.json

# login to aws ecr
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REGISTRY

# build, tag, and push app image to aws ecr
echo "docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG ./backend"
docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG ./backend
docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG

IMAGE_NAME=$ECR_REGISTRY/$ECR_REPOSITORY

# setup kubeconfig
aws eks --region $AWS_REGION update-kubeconfig --name $CLUSTER_NAME

# deploy backend server
helm upgrade --install rest-server ./deployments/helm-charts/backend \
  --set image.tag=$IMAGE_TAG \
  --set image.prod_repository=$IMAGE_NAME
