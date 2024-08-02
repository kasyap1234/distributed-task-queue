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
		return nil,err ; 


	}

	if err ==redis.Nil {
		return nil, err ; 
	}
	var job models.Job
	if err :=json.Unmarshal([]byte(result),&job); err!=nil {
		return nil, err ;
	}
	return &job,nil ; 

}
func (q *RedisQueue) EnqueueDeadLetter(job *models.Job)error{
	jobData, err :=json.Marshal(job); 
	if err !=nil {
		log.Fatal(err); 
		return err; 
	}
	return q.client.LPush(context.Background(),"dead_letter_queue",jobData).Err(); 

}
func (q *RedisQueue) GetJobs() ([]models.Job,error){
	result,err :=q.client.LRange(context.Background(),"job_queue",0,-1).Result()
	if err !=nil {
		return nil,err ;
	}
	jobs :=make([]models.Job,len(result))
	for i, result :=range result {
		var job models.Job 
		if err :=json.Unmarshal([]byte(result),&job); err!=nil {
			return nil, err ;
		}
		jobs[i] = job ;
		

	}
	return jobs,nil ; 
	
}
