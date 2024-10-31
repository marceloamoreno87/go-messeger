# Etapa de build: usa a imagem oficial do Golang baseada em Alpine Linux
FROM golang:alpine AS build

# Define um argumento de build chamado SERVICE_PATH
ARG SERVICE_PATH

# Define uma variável de ambiente para armazenar o caminho do serviço
ENV SERVICE_PATH_ENV=${SERVICE_PATH}

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos go.mod e go.sum do serviço para o diretório de trabalho
COPY ./go.mod ./go.sum ./

# Baixa as dependências do Go, utilizando cache para otimizar o processo
RUN --mount=type=cache,target=/go/shared/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

# Etapa de desenvolvimento: adiciona ferramentas úteis para desenvolvimento
FROM build AS development

# Instala o make, uma ferramenta de automação de build
RUN apk add --no-cache make

# Instala ferramentas adicionais para desenvolvimento, como air, delve, migrate e swag
RUN go install github.com/air-verse/air@latest && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Define o diretório de trabalho para o serviço específico
WORKDIR /app/$SERVICE_PATH_ENV

# Define o comando padrão para iniciar o contêiner em modo de desenvolvimento
CMD ["air", "-c", ".air.toml"]

# Etapa de build para produção: prepara o binário para produção
FROM build AS build-production

# Define um argumento de build chamado SERVICE_PATH
ARG SERVICE_PATH

ENV SERVICE_PATH_ENV=${SERVICE_PATH}

# Instala pacotes necessários, como shadow (para gerenciamento de usuários) e tzdata (para informações de fuso horário)
RUN apk add --no-cache shadow tzdata && \
    adduser -D -u 1001 appuser

# Copia todos os arquivos do contexto de build para o contêiner
COPY . .

# Compila o binário Go com otimizações para produção
RUN go build -x -tags netgo -o /main $SERVICE_PATH_ENV/main.go

# Etapa de produção: usa a imagem scratch para um contêiner mínimo
FROM scratch AS production

# Define um argumento de build chamado SERVICE_PATH
ARG SERVICE_PATH

ENV SERVICE_PATH_ENV=${SERVICE_PATH}

# Define o diretório de trabalho dentro do contêiner
WORKDIR /

# Copia o arquivo de senhas do sistema do contêiner de build para o contêiner de produção
COPY --from=build-production /etc/passwd /etc/passwd

# Copia o binário compilado para o contêiner de produção
COPY --from=build-production ./main .

# Copia as informações de fuso horário para o contêiner de produção
COPY --from=build-production /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime
COPY --from=build-production /usr/share/zoneinfo /usr/share/zoneinfo
# Copy Migration
COPY --from=build-production /app/$SERVICE_PATH_ENV/migrations /migrations

# Copia os certificados SSL para o contêiner de produção
COPY --from=build etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Define o usuário que executará o contêiner
USER appuser

# Define o comando padrão para iniciar o contêiner em produção
CMD ["/main"]