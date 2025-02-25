Banking Application

A simple banking application built in Golang with MongoDB, Postgres, Kafka, and Zookeeper for managing accounts and transactions.

🚀 Features

Account Creation: Submit account details to create a new user account.

Transactions: Record credit, debit, and transfer transactions.

Transaction History: Fetch the transaction history for a given account number.

Kafka Integration: Asynchronous processing of account and transaction requests via Kafka.

Database Integration: Persistent storage and retrieval of transaction data.

Swagger API documentation

Docker support for easy deployment

🛠️ Tech Stack

Golang

MongoDB, Postgres, Kafka, Zookeeper

Gorilla Mux (for routing)

Swagger (for API documentation)

Docker (for containerization)

🔧 Setup & Installation

1️⃣ Clone the Repository

git clone https://github.com/rockyisawesome/tempBankingApplication.git
cd tempBankingApplication

2️⃣ Install Dependencies

Install Docker Desktop to manage containers easily.

3️⃣ Set Up Environment Variables

No environment variables are required yet.

MongoDB

MONGO_URI=mongodb://admin:abcd@mongo:27017/ledger?authSource=admin
PORT=9091

🏗️ System Architecture

The application consists of four microservices that communicate asynchronously using Kafka:

***app.eraser.io link: https://app.eraser.io/workspace/3T8khE8tMTZelanhb53p?elements=yj82FvDw9TqVdk2Z2w8d9w

1️⃣ Account Producer (API Gateway)

Built with Gorilla Mux.

Handles incoming API requests.

Publishes requests to Kafka topics "account-creation" and "transactions" for further processing.

2️⃣ Account Service

Consumes messages from the "account-creation" Kafka topic.

Creates new user accounts and stores them in the database.

3️⃣ Transaction Service

Consumes messages from the "transactions" Kafka topic.

Processes transactions (credit, debit, transfers).

Publishes processed transactions to the "transaction-ledger" Kafka topic.

4️⃣ Ledger Service

Consumes messages from the "transaction-ledger" Kafka topic.

Creates transaction snapshots and saves them in MongoDB.

Maintains a transaction history log.

🐳 Run the Application

To start all services using Docker:

docker-compose up

This will spin up all required services in Docker containers, eliminating the need for manual dependencies.

Note: Check for the ledgerservice container, as it may require additional configuration to start correctly.

The server should be available at:

http://localhost:9091

📖 API Documentation

Swagger is available at:

http://localhost:9091/swagger/index.html

To regenerate Swagger documentation:

swag init

🔍 Testing

Run unit tests with:

go test -v ./...

(Currently working on test cases)

📜 License

This project is licensed under the MIT License.

🤝 Contributing

Feel free to fork the repo, open issues, and submit PRs to improve this project!

🔗 Connect

📧 Email: pandeychandransh@gmail.com🐙 GitHub: rockyisawesome📞 Mobile: 8468950657
