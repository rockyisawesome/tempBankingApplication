Temp Banking Application

A simple banking application built in Golang with MongoDB for managing accounts and transactions.

🚀 Features

User account creation

Deposit & withdrawal functionality

Transaction history tracking

Swagger API documentation

Docker support for easy deployment

🛠️ Tech Stack

Golang

MongoDB

Gorilla Mux (for routing)

Swagger (for API documentation)

Docker (for containerization)

🔧 Setup & Installation

1️⃣ Clone the Repository

git clone https://github.com/rockyisawesome/tempBankingApplication.git
cd tempBankingApplication

2️⃣ Install Dependencies

go mod tidy

3️⃣ Set Up Environment Variables

Create a .env file in the project root and add:

MONGO_URI=mongodb://admin:abcd@mongo:27017/ledger?authSource=admin
PORT=9091

4️⃣ Run the Application

go run main.go

The server should start on http://localhost:9091

🐳 Running with Docker

1️⃣ Build Docker Image

docker build -t banking-app .

2️⃣ Run the Container

docker run -p 9091:9091 banking-app

📖 API Documentation

Swagger is available at:

http://localhost:9091/swagger/index.html

To regenerate Swagger docs, run:

swag init

🔍 Testing

Run unit tests with:

go test -v ./...

📜 License

This project is licensed under the MIT License.

🤝 Contributing

Feel free to fork the repo, open issues, and submit PRs to improve this project!

🔗 Connect

📧 Email: your-email@example.com🐙 GitHub: rockyisawesome
