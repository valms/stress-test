# Stress Test CLI

Uma ferramenta de linha de comando em Go para realizar testes de carga em serviços web, desenvolvida como parte do curso
Pós Go Expert - 2024 da Full Cycle.

## 🎯 Objetivo

Realizar testes de carga em serviços web através de uma interface CLI, permitindo configurar:

- URL do serviço alvo
- Número total de requisições
- Nível de concorrência (threads simultâneas)

## ✨ Funcionalidades

- Execução de requisições HTTP concorrentes
- Controle preciso do número de threads
- Coleta de métricas de performance
- Relatório detalhado de execução
- Suporte a HTTPS
- Execução via Docker

## 🚀 Como Usar

### Via Docker

```bash
# Build da imagem
docker build -t stress-test .

# Executar teste
docker run stress-test \
    --url=http://example.com \
    --requests=1000 \
    --threads=10
```

### Localmente

```bash
# Build
go build -o stress-test .

# Executar
./stress-test \
    --url=http://example.com \
    --requests=1000 \
    --threads=10
```

## 📊 Exemplo de Saída

```
🚀 Iniciando Teste de Performance
================================
URL: http://example.com
Total de Requisições: 1000
Threads Concorrentes: 10

📊 Relatório de Performance
==========================

⏱️  Duração Total: 5.234s
📈 Requisições Executadas: 1000

✅ Distribuição por Status HTTP:
   Status 200: 950 requisições
   Status 404: 30 requisições
   Status 500: 20 requisições

📈 Resumo Final:
   • Requisições bem-sucedidas (2xx): 950
   • Total de falhas: 50
   • Taxa de sucesso: 95.00%
```

## 🛠 Tecnologias Utilizadas

- Go 1.22
- Docker
- Concorrência nativa do Go
- HTTP Client otimizado

## 📁 Estrutura do Projeto

```
.
├── cmd
│   └── main.go
├── service
│   └── benchmark.go
├── Dockerfile
├── go.mod
└── README.md
```

## ⚙️ Configurações

| Flag       | Descrição                      | Obrigatório |
|------------|--------------------------------|-------------|
| --url      | URL do serviço a ser testado   | Sim         |
| --requests | Número total de requisições    | Sim         |
| --threads  | Número de threads concorrentes | Sim         |

## 📝 Licença

Este projeto é parte do curso Pós Go Expert - 2024 da Full Cycle.