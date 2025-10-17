# My Kafka Microservices

This project consists of a set of microservices that communicate through Kafka. It includes two services, Service 1 and Service 2, which are deployed on a Kubernetes cluster. 

## Project Structure

```
my-kafka-microservices
├── services
│   ├── service1
│   │   ├── main.go
│   │   └── config.yaml
│   └── service2
│       ├── main.go
│       └── config.yaml
├── k8s
│   ├── service1-deployment.yaml
│   ├── service2-deployment.yaml
│   └── kafka-deployment.yaml
└── README.md
```

## Services Overview

### Service 1
- **Description**: Service 1 sets up an HTTP server that listens for incoming requests, processes them, and sends messages to a Kafka topic.
- **Configuration**: The configuration settings for Service 1, including Kafka broker addresses and topic names, are defined in `services/service1/config.yaml`.

### Service 2
- **Description**: Service 2 acts as a Kafka consumer that listens for messages from a specified topic and processes them accordingly.
- **Configuration**: The configuration settings for Service 2, including Kafka broker addresses and topic names, are defined in `services/service2/config.yaml`.

## Kubernetes Deployment

The services and Kafka are deployed on a Kubernetes cluster using the following deployment configurations:

- **Service 1 Deployment**: Defined in `k8s/service1-deployment.yaml`.
- **Service 2 Deployment**: Defined in `k8s/service2-deployment.yaml`.
- **Kafka Deployment**: Defined in `k8s/kafka-deployment.yaml`.

## Setup Instructions

1. **Clone the repository**:
   ```
   git clone <repository-url>
   cd my-kafka-microservices
   ```

2. **Deploy Kafka**:
   Apply the Kafka deployment configuration:
   ```
   kubectl apply -f k8s/kafka-deployment.yaml
   ```

3. **Deploy Service 1**:
   Apply the Service 1 deployment configuration:
   ```
   kubectl apply -f k8s/service1-deployment.yaml
   ```

4. **Deploy Service 2**:
   Apply the Service 2 deployment configuration:
   ```
   kubectl apply -f k8s/service2-deployment.yaml
   ```

## Testing Communication

To test the communication between the services through Kafka, you can send requests to Service 1 and observe the messages being processed by Service 2. 

Ensure that you have the necessary tools installed (e.g., `kubectl`, Kafka client) to interact with the Kubernetes cluster and Kafka. 

## License

This project is licensed under the MIT License.