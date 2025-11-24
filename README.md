# 🧩 Ecom-Micro — Microservices-Based E-Commerce Platform (Go + RabbitMQ + Docker)

Ecom-Micro is a fully containerized microservices architecture built using **Go**, **RabbitMQ**, and **Docker Compose**.  
Each domain (Users, Products, Orders, Payments, Notifications) is implemented as an independent service with its own database and communication layer.  
The project is designed for learning and understanding real-world microservices patterns such as:

- Service-level isolation  
- Event-driven communication  
- JWT-based authentication  
- Domain-based responsibility separation  
- CI/CD-friendly structure  

---

## 🚀 Architecture Overview

The system follows a distributed microservices architecture consisting of:

### **1. User Service**
- Handles registration & login  
- Password hashing with bcrypt  
- JWT authentication  
- Issues access tokens for other services  

### **2. Product Service**
- Manages product catalog  
- CRUD operations  
- Dedicated Postgres database  

### **3. Order Service**
- Places & manages orders  
- Validates tokens with User Service  
- Publishes events to RabbitMQ  

### **4. Payment Service**
- Processes payments  
- Subscribes to order events  
- Publishes payment confirmations  

### **5. Notification Service**
- Subscribes to order & payment events  
- Sends email or in-app notifications  

### **6. RabbitMQ Message Broker**
- Internal messaging system  
- Event-driven communication backbone  

### **7. Docker Compose**
- Spins up all services and infrastructure  
- Local development environment for microservices  

---

## 🧱 Technologies Used

| Component | Technology |
|----------|------------|
| Language | Go (Golang) |
| Messaging | RabbitMQ |
| Databases | Postgres (one DB per service) |
| Containerization | Docker & Docker Compose |
| Auth | JWT + bcrypt |
| Communication | REST + AMQP Events |
| Architecture | Microservices + Event-Driven + DDD |

---

## 📂 Folder Structure

ecom-micro/
│
├── user-service/
├── product-service/
├── order-service/
├── payment-service/
├── notify-service/
│
├── docker-compose.yml
└── .gitignore


Each service contains its own internal modules such as:

- controllers/
- routes/
- models/
- config/
- utils/
- middleware/
- Dockerfile

---

## 🛠️ Running the Project Locally

### Prerequisites

- Docker Desktop
- Docker Compose
- Go (optional, only for manual service runs)


### ▶️ Start All Services

Run this command in the project root:

docker-compose up --build

This will start:

All microservices
All Postgres databases
RabbitMQ
pgadmin

### RabbitMQ Dashboard

- http://localhost:15672
- username: guest
- password: guest


### pgAdmin Dashboard

- http://localhost:5050
- email: admin@example.com
- password: admin


## 🔑 Authentication Flow

User registers via /auth/register
User logs in via /auth/login
User receives a JWT:

Authorization: Bearer <token>

Other services validate the token through the User Service

Validated requests can then proceed (place orders, view products, etc.)


## 📬 Event-Driven Flow Example (Order → Payment → Notification)

Order Service publishes order.created

Payment Service consumes the event & processes payment

Payment Service publishes payment.completed

Notification Service listens and notifies the user

This ensures loose coupling and scalability.


## 🧪 Running Individual Services

cd user-service
go run main.go

Rebuild a service
docker-compose build user-service

Restart a service
docker-compose up user-service


## 📘 Future Enhancements

API Gateway (Kong / Traefik)
Kubernetes deployment
Distributed tracing with Jaeger
Elasticsearch + Kibana for logs
Consul for service discovery


## 🤝 Contributing

Pull requests are welcome.
Before making major changes, open an issue for discussion.


## 📄 License

This project is licensed under the MIT License.


## ✨ Author

Naveen Hettiwaththa
Software Engineer | Undergraduate Computer Science Student (Final Year) | Go Enthusiast