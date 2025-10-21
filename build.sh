#!/bin/bash
set -e

# Default values
REGISTRY="${DOCKER_REGISTRY:-docker.io}"
USERNAME="${DOCKER_USERNAME:-yourusername}"
TAG="${IMAGE_TAG:-latest}"
NS="go-kafka-lab"

SVC1_IMG="${REGISTRY}/${USERNAME}/goservice1:${TAG}"
SVC2_IMG="${REGISTRY}/${USERNAME}/goservice2:${TAG}"

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTION]"
    echo "Options:"
    echo "  --full, -f       Full build and deploy (build images, push, and deploy to K8s)"
    echo "  --deploy, -d     Deploy only to K8s (skip building and pushing images)"
    echo "  --help, -h       Show this help message"
    echo ""
    echo "Environment variables:"
    echo "  DOCKER_REGISTRY  Docker registry (default: docker.io)"
    echo "  DOCKER_USERNAME  Docker username (default: yourusername)"
    echo "  IMAGE_TAG        Image tag (default: latest)"
}

# Function to build and push images
build_and_push() {
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
}

# Function to deploy to K8s
deploy_to_k8s() {
    echo "=== Deploying to Kubernetes ==="
    kubectl apply -f k8s/namespace.yaml
    kubectl apply -f k8s/kafka-deployment.yaml
    kubectl wait --for=condition=ready pod -l app=kafka -n ${NS} --timeout=180s || true
    kubectl apply -f k8s/service1-deployment.yaml
    kubectl apply -f k8s/service1-service.yaml
    kubectl apply -f k8s/service2-deployment.yaml

    echo "=== Checking deployment status ==="
    kubectl get pods -n ${NS}
    kubectl get services -n ${NS}
}

# Parse command line arguments
case "${1:-}" in
    --full|-f)
        echo "=== Full Build and Deploy ==="
        build_and_push
        deploy_to_k8s
        ;;
    --deploy|-d)
        echo "=== Deploy Only ==="
        deploy_to_k8s
        ;;
    --help|-h)
        show_usage
        exit 0
        ;;
    "")
        # Default behavior (backward compatibility) - full build and deploy
        echo "=== Full Build and Deploy (default) ==="
        build_and_push
        deploy_to_k8s
        ;;
    *)
        echo "Error: Unknown option '$1'"
        show_usage
        exit 1
        ;;
esac

echo "=== Done ==="
echo "Service1 external endpoint: kubectl get service service1-external -n ${NS}"
echo "Test with: curl -X POST -H 'Content-Type: application/json' -d '{\"content\":\"Test message\"}' http://<EXTERNAL-IP>:30080/send"
echo "Or port-forward: kubectl port-forward svc/service1 8080:8080 -n ${NS}"
echo "Then: curl -X POST http://localhost:8080/send -d '{\"content\":\"test\"}' -H 'Content-Type: application/json'"