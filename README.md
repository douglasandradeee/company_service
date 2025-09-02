# 🏢 Company Service API

API RESTful para gerenciamento de empresas com validação de CNPJ.

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Pré-requisitos](#pré-requisitos)
- [Instalação](#instalação)
- [Executando o Projeto](#executando-o-projeto)
- [Endpoints](#endpoints)
- [Testes](#testes)
- [Configuração](#configuração)

## 🌐 Visão Geral

Este projeto é uma API RESTful desenvolvida para gerenciar informações de empresas, com foco em validação de CNPJ (Cadastro Nacional da Pessoa Jurídica). A API permite operações CRUD (Criar, Ler, Atualizar, Deletar) para entidades de empresas.

## 🛠️ Pré-requisitos

Antes de começar, certifique-se de ter instalado:

- [Go](https://golang.org/dl/) (versão 1.20 ou superior)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## 📦 Instalação

### Clonar o Repositório

```bash
git clone https://github.com/douglasandradeee/company-service.git
cd company-service
```

## 🚀 Executando o Projeto

### Método 1: Docker Compose (Recomendado)

```bash
# Inicia todos os serviços definidos no docker-compose.yml
docker-compose up -d

# Para parar os serviços
docker-compose down
```

### Método 2: Execução Local

```bash
# Baixar dependências
go mod download

# Compilar o projeto
go build -o company-service

# Executar
./company-service
```

## 🌍 Endpoints

### Host da API

- **URL Base**: `http://localhost:8080`
- **Porta Padrão**: 8080
- **Ambiente**: Desenvolvimento

### Empresas

- `GET /companies`: Listar todas as empresas
- `GET /companies/{id}`: Buscar empresa por ID
- `POST /companies`: Criar nova empresa
- `PUT /companies/{id}`: Atualizar empresa existente
- `DELETE /companies/{id}`: Remover empresa

### Saúde

- `GET /health`: Verificar status da aplicação

## 🐰 Integração com RabbitMQ

Esta API está integrada com RabbitMQ para processamento de mensagens assíncronas. A cada operação CRUD (Criar, Ler, Atualizar, Deletar) realizada na API, uma mensagem é enviada para uma fila do RabbitMQ.

### Fluxo de Mensagens

- **Criação de Empresa**: Envia mensagem para a fila `company.created`
- **Atualização de Empresa**: Envia mensagem para a fila `company.updated`
- **Exclusão de Empresa**: Envia mensagem para a fila `company.deleted`

### Configurações de RabbitMQ

As seguintes variáveis de ambiente são usadas para configurar a conexão:

- `RABBITMQ_HOST`: Endereço do servidor RabbitMQ
- `RABBITMQ_PORT`: Porta do RabbitMQ (padrão: 5672)
- `RABBITMQ_USER`: Usuário para autenticação
- `RABBITMQ_PASSWORD`: Senha para autenticação
- `RABBITMQ_VHOST`: Virtual Host do RabbitMQ (opcional)

## 🧪 Testes

Para executar os testes:

```bash
go test ./...
```

## 🔧 Configuração

### Variáveis Disponíveis

- `SERVER_PORT`: Porta do servidor (padrão: 8080)
- `MONGO_URI`: URI de conexão com MongoDB (padrão: mongodb://localhost:27017)
- `MONGO_DB`: Nome do banco de dados MongoDB (padrão: company_db)
- `MONGO_COLLECTION`: Nome da coleção de empresas (padrão: companies)
- `RABBITMQ_URI`: URI de conexão com RabbitMQ (padrão: amqp://guest:guest@localhost:5672/)
- `QUEUE_NAME`: Nome da fila de eventos (padrão: company_events)
- `LOG_LEVEL`: Nível de log da aplicação (padrão: info)
- `SHUTDOWN_TIMEOUT`: Tempo limite para desligamento (padrão: 10s)
- `READ_TIMEOUT`: Tempo limite de leitura (padrão: 5s)
- `WRITE_TIMEOUT`: Tempo limite de escrita (padrão: 10s)
- `IDLE_TIMEOUT`: Tempo limite ocioso (padrão: 60s)

### Precedência das Configurações

1. Variáveis de ambiente (maior prioridade)
2. Arquivo `.env`
3. Valores padrão definidos no código

```

### Notas Importantes

- As configurações padrão são adequadas para desenvolvimento local
- Variáveis de ambiente têm precedência sobre o arquivo `.env`
- Ajuste os timeouts conforme necessidade da sua aplicação
- Use credenciais seguras em ambientes de produção
```
