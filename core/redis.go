package core

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Configuração básica do Redis
type RedisConfig struct {
	DSN string
}

// Estrutura para implementar o cliente Redis e manter a conexão e configuração.
type RedisClient struct {
	config *RedisConfig
	client *redis.Client
	ctx    context.Context
}

// Implementação do método Connect, aplicando o princípio de responsabilidade única.
func (r *RedisClient) Connect() error {
	r.ctx = context.Background()
	r.client = redis.NewClient(&redis.Options{
		Addr: r.config.DSN,
	})

	// Testa a conexão com o Redis
	_, err := r.client.Ping(r.ctx).Result()
	if err != nil {
		return fmt.Errorf("não foi possível conectar ao Redis: %w", err)
	}
	fmt.Println("Conectado ao Redis com sucesso")
	return nil
}

// Implementação do método Set para armazenar dados em cache.
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(r.ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("erro ao definir valor no Redis: %w", err)
	}
	return nil
}

// Implementação do método Get para recuperar dados do cache.
func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("chave não encontrada")
	} else if err != nil {
		return "", fmt.Errorf("erro ao obter valor do Redis: %w", err)
	}
	return val, nil
}

// Implementação do método Close para encerrar a conexão.
func (r *RedisClient) Close() error {
	return r.client.Close()
}
