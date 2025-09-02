🏢 Company Service API

📋 Sobre o Projeto

API RESTful completa para gerenciamento de empresas com validação rigorosa de CNPJ e sistema de eventos em tempo real.

🏗️ Arquitetura

📦 company-service

├── 🌐 API Layer (HTTP Handlers)
├── 🔧 Service Layer (Lógica de Negócio + Eventos)
├── 🗄️ Repository Layer (MongoDB)
├── 📨 Messaging Layer (RabbitMQ)
└── 🧪 Test Suite

# Clone e execução em 1 comando

git clone [def]
cd company-service
docker-compose up -d

# Verificar saúde da aplicação

curl http://localhost:8080/health

📊 Endpoints

Método Endpoint Descrição Status Code
POST /companies Criar empresa 201 Created
GET /companies/{id} Buscar empresa 200 OK
PUT /companies/{id} Atualizar empresa 200 OK
DELETE /companies/{id} Deletar empresa 204 No Content
GET /companies Listar empresas 200 OK

🎯 Funcionalidades Principais

✅ Validação de CNPJ

Formatação e limpeza automática

Impedimento de CNPJs duplicados

📨 Sistema de Eventos

Cada operação gera evento no RabbitMQ:

Cadastro de EMPRESA {NOME} - Para novas empresas

Edição da EMPRESA {NOME} - Para atualizações

Exclusão da EMPRESA {NOME} - Para remoções

🧪 Suite de Testes

# Executar testes

go test ./... -v

# Cobertura de testes

go test ./... -coverprofile=coverage.out

go tool cover -html=coverage.out

Cobertura atual: 80%+ (Domain, Repository, Service Layers)

🛠️ Tecnologias Utilizadas

Camada Tecnologia Finalidade

Linguagem Go 1.24 Performance e concorrência

Banco MongoDB Persistência de documentos

Mensageria RabbitMQ Eventos assíncronos

Container Docker Ambiente isolado

Orquestração Docker Compose Multi-container setup

Logging Zap Logs estruturados e performáticos

🔧 Configuração

Variáveis de Ambiente (/.env)

env

SERVER_PORT=8080

MONGO_URI=mongodb://mongodb:27017

MONGO_DB=company_db

MONGO_COLLECTION=companies

RABBITMQ_URI=amqp://rabbitmq:5672/

QUEUE_NAME=company_events

LOG_LEVEL=info

📦 Estrutura do Projeto

company-service/
├── cmd/api/main.go # Entrypoint da aplicação
├── internal/
│ ├── domain/ # Entidades e regras de negócio
│ ├── repository/ # Persistência (MongoDB)
│ ├── service/ # Lógica de aplicação
│ ├── handler/ # HTTP Handlers
│ ├── messaging/ # Interface de mensageria
│ └── server/ # Configuração do servidor
├── pkg/utils/ # Utilitários (validação CNPJ)
└── tests/ # Testes unitários e integração

# Health check da aplicação

curl http://localhost:8080/health

# Health check MongoDB

docker-compose exec mongodb mongosh --eval "db.adminCommand('ping')"

# Health check RabbitMQ

curl -u guest:guest http://localhost:15672/api/healthchecks/node

[def]: https://github.com/douglasandradeee/company_service
