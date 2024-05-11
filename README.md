# Microservices - E-commerce Order Processing

## Components:

1. **Microservices Architecture**:

    - Multiple microservices: User Service, Order Service, Product Service, and Inventory Service. This architecture allows for scalability, flexibility, and easier management.

2. **Fiber for HTTP API Gateway**:

    - Fiber to create a lightweight, high-performance HTTP API gateway that acts as a central entry point for clients (e.g., frontend applications, mobile apps). This API gateway handles authentication, authorization, and routing requests to the appropriate microservices.

3. **PostgreSQL & GORM for Data Storage**:

    - Each microservice has its own PostgreSQL database for data storage, using GORM as the ORM layer to interact with the database. This decouples the microservices and provides data isolation and consistency.

4. **Kafka for Event-Driven Communication**:

    - Kafka for asynchronous communication between microservices, allowing for decoupled services that can publish and consume events, enabling real-time processing and scalability.

5. **gRPC for Inter-Service Communication**:

    - gRPC for synchronous communication between microservices when direct communication is required, such as getting product details from the Product Service when placing an order.

## Implementation Details:

1. **User Service**:

    - This microservice handles user registration, authentication, and authorization. It interacts with PostgreSQL through GORM to manage user data.
    - Provides REST endpoints for user-related actions via Fiber.
    - Communicates with other microservices using gRPC and Kafka when user-related events occur (e.g., "User Registered").

2. **Product Service**:

    - Manages product data, including product details, inventory, and pricing.
    - Exposes REST endpoints for product-related operations via Fiber.
    - Communicates with other services using Kafka for product-related events and gRPC for real-time data requests.

3. **Order Service**:

    - Handles order placement, tracking, and status updates.
    - Uses Kafka to publish "Order Placed" events and listens to Kafka topics for inventory changes.
    - Communicates with Product Service via gRPC to fetch product information for order processing.

4. **Inventory Service**:
    - Manages inventory levels and adjusts stock based on orders.
    - Listens to Kafka topics for order-related events and publishes "Inventory Changed" events.
    - Provides REST endpoints for inventory management via Fiber.

---

-   **Authentication and Authorization**: JWT-based authentication and authorization to ensure secure access to the API gateway and microservices.

-   **Logging and Monitoring**: Prometheus and Grafana to track the behavior of microservices and Kafka events.

<!-- export PATH="$PATH:$(go env GOPATH)/bin" -->
<!-- protoc --go_out=. --go-grpc_out=. proto/user.proto -->
