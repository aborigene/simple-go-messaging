#!/bin/bash
set -e

REGISTRY="${DOCKER_REGISTRY:-docker.io}"
USERNAME="${DOCKER_USERNAME:-yourusername}"
TAG="${IMAGE_TAG:-latest}"
NS="kafka-demo"

SVC1_IMG="${REGISTRY}/${USERNAME}/service1:${TAG}"
SVC2_IMG="${REGISTRY}/${USERNAME}/service2:${TAG}"

echo "=== Building Service1 ==="
cd services/service1
docker build -t ${SVC1_IMG} .
cd ../..

echo "=== Building Service2 ==="
cd services/service2
docker build -t ${SVC2_IMG} .
cd ../..

echo "=== Pushing Images ==="
docker push ${SVC1_IMG}
docker push ${SVC2_IMG}

echo "=== Updating Deployments ==="
sed -i.bak "s|image:.*service1.*|image: ${SVC1_IMG}|g" k8s/service1-deployment.yaml
sed -i.bak "s|image:.*service2.*|image: ${SVC2_IMG}|g" k8s/service2-deployment.yaml
rm -f k8s/*.bak

echo "=== Deploying to Kubernetes ==="
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/kafka-deployment.yaml
kubectl wait --for=condition=ready pod -l app=kafka -n ${NS} --timeout=180s || true
kubectl apply -f k8s/service1-deployment.yaml
kubectl apply -f k8s/service2-deployment.yaml

echo "=== Done ==="
echo "Test: kubectl port-forward svc/service1 8080:8080 -n ${NS}"
echo "Then: curl -X POST http://localhost:8080/messages -d '{\"content\":\"test\"}' -H 'Content-Type: application/json'"