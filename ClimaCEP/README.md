# Sistema de Temperatura por CEP com OpenTelemetry e Zipkin

Este projeto demonstra um sistema de microserviços em Go para obter a temperatura de uma cidade a partir de um CEP, utilizando OpenTelemetry para tracing distribuído e Zipkin para visualização.

## Arquitetura

O sistema é composto por dois serviços principais:

- **Serviço A (servicoa):** Atua como o ponto de entrada. Ele recebe uma requisição POST com um CEP, valida o formato do CEP e o encaminha para o Serviço B.
- **Serviço B (servicob):** É responsável pela orquestração. Ele recebe o CEP, consulta a API ViaCEP para obter o nome da cidade e, em seguida, utiliza a WeatherAPI para buscar a temperatura atual da cidade. Por fim, calcula as temperaturas em Celsius, Fahrenheit e Kelvin e retorna a resposta.

## Requisitos

- Docker
- Docker Compose
- Uma chave da WeatherAPI (você pode obtê-la em [https://www.weatherapi.com/](https://www.weatherapi.com/))

## Configuração

1.  **Obtenha sua chave da WeatherAPI:**
    Visite [https://www.weatherapi.com/](https://www.weatherapi.com/) e crie uma conta para obter sua API Key.
    Substituir WEATHER_API_KEY diretamente no `docker-compose.yml` pela sua chave.

## Como Rodar o Projeto

2.  **Construa e inicie os contêineres Docker:**
    No diretório raiz do projeto, execute:
    ```bash
    docker-compose up --build
    ```
    Isso irá construir as imagens Docker para o Serviço A e Serviço B, e iniciar os contêineres, incluindo o Zipkin.

    Aguarde até que os serviços estejam totalmente inicializados. Você verá logs indicando que "Serviço A rodando na porta 8082" e "Serviço B rodando na porta 8081".

## Testando a Aplicação

Você pode usar `curl` ou um cliente HTTP como Postman/Insomnia para testar os endpoints.

O endpoint principal para interagir é o do **Serviço A**.

### Cenários de Teste:

#### 1. Sucesso (CEP válido e encontrado)

**Requisição (POST para `http://localhost:8080/cep`):**

```json
{
    "cep": "29902555"
}