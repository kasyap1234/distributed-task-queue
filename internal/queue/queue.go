package queue

import (
	"context"
	"distributed-task-queue/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
	queueName string 
}
var ErrQueueEmpty = errors.New("queue is empty")

func NewRedisQueue(addr string,queueName string) *RedisQueue {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisQueue{client: client,queueName: queueName}


}

func (q *RedisQueue) Enqueue(ctx context.Context,jobType string,payload interface {},opts ...JobOption)(*models.Job,error) {
job := &models.Job {
	ID : uuid.New().String(), 
	Type: jobType, 
	Status : "pending", 
	CreatedAt : time.Now(), 
	MaxRetries : 3 , 
	Priority : 1, 
}
for _,opt := range opts {
	opt(job); 
}
payloadBytes ,err :=json.Marshal(payload);
if err !=nil {
	return nil,err; 
}
job.Payload=payloadBytes; 
jobBytes,err :=json.Marshal(job); 
if err !=nil {
	return nil, err ; 
}
score :=float64(time.Now().Add(time.Until(job.DelayUntil)).UnixNano())
 q.client.ZAdd(context.Background(),q.queueName,redis.Z{
		Score: score, 
		Member: jobBytes,
	}).Err()
	return 
}
func (q *RedisQueue) Dequeue()(*models.Job,error){
	result,err := q.client.ZRangeByScore(context.Background(),"job_queue",&redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprint(time.Now().UnixNano()),
		Count : 1 , 
	}).Result()
	if err == redis.Nil {
		return nil,ErrQueueEmpty
	}
	if len(result) == 0 {
		return nil,ErrQueueEmpty
	}
	var job models.Job 
	if err :=json.Unmarshal([]byte(result[0]),&job); err !=nil {
		return nil,err
	}
	return &job,nil
}
func (q *RedisQueue) EnqueueDeadLetter(job *models.Job)error {
	jobData, err :=json.Marshal(job); 
	if err !=nil {
		return err; 
	}
return q.client.LPush(context.Background(),"dead_letter_queue",jobData).Err()

}
func (q *RedisQueue) GetJobs() ([]models.Job, error){
	jobsData,err := q.client.ZRange(context.Background(),"job_queue",0,-1).Result(); 
	if err !=nil {
		return nil , err 
	}
	jobs :=make([]models.Job,len(jobsData)); 
	for i,jobData :=range jobsData {
		var job models.Job 
		if err :=json.Unmarshal([]byte(jobData),&job); err !=nil {
			return nil,err
		}
		jobs[i] = job

	}
	return jobs,nil
}