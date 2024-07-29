package queue 
import (
	"context"
	"encoding/json"
	"log"
	"github.com/go-redis/redis/v8"
	"distributed-task-queue/internal/models"
)
type RedisQueue struct {
	client *redis.Client 

}
func NewRedisQueue()