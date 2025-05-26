# Rate Limiter em Go com Redis

Este projeto implementa um **Rate Limiter (limitador de requisições)** em Go que pode ser configurado para restringir o número de requisições por **IP** ou por **Token de Acesso**. Ele funciona como um middleware para servidores web e utiliza o Redis como mecanismo de armazenamento.

---

## Funcionalidades

- Limitação de requisições por endereço IP
- Limitação de requisições por token de acesso (`API_KEY`)
- Configuração de limites e tempo de bloqueio via `.env`
- Redis como mecanismo de persistência
- Arquitetura desacoplada por estratégia
- Middleware reutilizável
- Suporte a Docker e Docker Compose

---

## Pré-requisitos

- [Go 1.23.3+](https://go.dev/dl/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## Configuração

Crie um arquivo `.env` na raiz do projeto com os seguintes parâmetros:

## .env
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
RATE_LIMIT=5
BLOCK_DURATION=300
TOKEN_RATE_LIMIT_abc123=10
TOKEN_BLOCK_DURATION_abc123=60

## Execução
Após executar o programa, abra um novo prompt e execute o comando curl -H "API_KEY: abc123" http://localhost:8081
