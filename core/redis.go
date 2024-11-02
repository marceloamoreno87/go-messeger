package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

/*
Definição de variáveis de erro específicas para operações com o Redis.
Essas variáveis são usadas para fornecer mensagens de erro detalhadas.
*/
var (
	ErrRedisURLParse    = errors.New("redis.url_parse_failed: erro ao analisar a URL do Redis")
	ErrRedisConnection  = errors.New("redis.connection_failed: erro ao conectar ao Redis")
	ErrRedisSetValue    = errors.New("redis.set_value_failed: erro ao definir valor no Redis")
	ErrRedisDeleteKey   = errors.New("redis.delete_key_failed: erro ao deletar chave no Redis")
	ErrRedisGetValue    = errors.New("redis.get_value_failed: erro ao obter valor do Redis")
	ErrRedisKeyNotFound = errors.New("redis.key_not_found: chave não encontrada")
)

/*
Estrutura para implementar o cliente Redis e manter a conexão e configuração.
*/
type RedisClient struct {
	DSN    string
	client *redis.Client
	ctx    context.Context
}

/*
Implementação do método Connect, aplicando o princípio de responsabilidade única.
Estabelece uma conexão com o Redis e retorna um ponteiro para o RedisClient e um erro, se houver.
*/
func (r *RedisClient) Connect() (*RedisClient, error) {
	r.ctx = context.Background()

	options, err := redis.ParseURL(r.DSN)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRedisURLParse, err)
	}

	r.client = redis.NewClient(options)

	// Testa a conexão com o Redis
	_, err = r.client.Ping(r.ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRedisConnection, err)
	}
	return r, nil
}

/*
Implementação do método Set para armazenar dados em cache.
Define um valor no Redis com uma chave específica e um tempo de expiração.
Retorna um erro, se houver.
*/
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(r.ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRedisSetValue, err)
	}
	return nil
}

/*
Implementação do método Del para deletar dados do cache.
Remove um valor do Redis com uma chave específica.
Retorna um erro, se houver.
*/
func (r *RedisClient) Del(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRedisDeleteKey, err)
	}
	return nil
}

/*
Implementação do método Get para recuperar dados do cache.
Obtém um valor do Redis com uma chave específica.
Retorna o valor e um erro, se houver.
*/
func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", ErrRedisKeyNotFound
	} else if err != nil {
		return "", fmt.Errorf("%w: %v", ErrRedisGetValue, err)
	}
	return val, nil
}

/*
Implementação do método Close para encerrar a conexão.
Fecha a conexão com o Redis.
Retorna um erro, se houver.
*/
func (r *RedisClient) Close() error {
	return r.client.Close()
}
