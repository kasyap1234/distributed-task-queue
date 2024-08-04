package queue

import (
	"context"
	"distributed-task-queue/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
}
var ErrQueueEmpty = errors.New("queue is empty")

func NewRedisQueue(addr string) *RedisQueue {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisQueue{client: client}

}

func (q *RedisQueue) Enqueue(job *models.Job) error {
	jobData,err :=json.Marshal(job); 
	if err !=nil {
		return err; 
	}
	score :=float64(time.Now().Add(time.Until(job.DelayUntil)).UnixNano())
    return q.client.ZAdd(context.Background(),"job_queue",redis.Z{
		Score: score, 
		Member: jobData, 
	}).Err()
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