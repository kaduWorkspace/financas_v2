# Use uma imagem base oficial do Golang para construir o projeto
FROM golang:1.20 AS builder

# Defina o diretório de trabalho dentro do container
WORKDIR /app

# Copie o código do projeto para o diretório de trabalho
COPY . .

# Baixe as dependências do projeto (módulos Go)
RUN go mod tidy

# Compile o projeto
RUN go build -o meu-aplicativo

# Segunda etapa: criar uma imagem final menor para execução
FROM debian:bullseye-slim

# Defina o diretório de trabalho da imagem final
WORKDIR /app

# Copie o binário compilado da etapa anterior
COPY --from=builder /app/meu-aplicativo .

# Expõe a porta onde o servidor vai rodar (substitua pela porta correta, ex: 8080)
EXPOSE 3000

# Comando para rodar o aplicativo
CMD ["./goravel"]

