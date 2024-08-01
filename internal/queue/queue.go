package queue 
import (
	"context"
	"encoding/json"
	"log"
	"github.com/redis/go-redis/v9"
	"distributed-task-queue/internal/models"
)
type RedisQueue struct {
	client *redis.Client 

}
func NewRedisQueue(addr string) *RedisQueue {
	client :=redis.NewClient(&redis.Options{
		Addr:addr ,

	})
	return &RedisQueue{client :client }; 


}

func (q *RedisQueue) Enqueue(job *models.Job) error {
	jobData,err :=json.Marshal(job); 
	if err !=nil {
		return err; 
	}
	return q.client.LPush(context.Background(),"job_queue",jobData).Err()


}
func (q *RedisQueue) Dequeue() (*models.Job,error){
	result,err :=q.client.RPop(context.Background(),"job_queue").Result()
	if err !=nil {
		return nil, 

	}
	if err ==redis.Nil {
		return nil, ErrQueueEmpty
	}
	var job models.Job
	if err :=json.Unmarshal([]byte(result),&job); err!=nil {
		return nil, err ;
	}
	return &job,nil ; 
	
}
