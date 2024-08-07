package handlers

import (
	"context"
	"distributed-task-queue/internal/models"
	"distributed-task-queue/internal/storage"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/smtp"
	"os"

	"github.com/nfnt/resize"
)

func HandleEmailJob(ctx context.Context, job *models.Job) error {
    var emailJob models.EmailJob
    err := json.Unmarshal(job.Payload, &emailJob)
    if err != nil {
        return fmt.Errorf("failed to unmarshal email job: %w", err)
    }

    auth := smtp.PlainAuth("", "your-email@example.com", "your-password", "smtp.example.com")
    to := []string{emailJob.To}
    msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", emailJob.To, emailJob.Subject, emailJob.Body))

    err = smtp.SendMail("smtp.example.com:587", auth, "your-email@example.com", to, msg)
    if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

    log.Printf("Email sent to %s", emailJob.To)
    return nil
}

func HandleImageResizeJob(ctx context.Context, job *models.Job) error {
    var imageJob models.ImageResizeJob
    err := json.Unmarshal(job.Payload, &imageJob)
    if err != nil {
        return fmt.Errorf("failed to unmarshal image resize job: %w", err)
    }

    file, err := os.Open(imageJob.InputPath)
    if err != nil {
        return fmt.Errorf("failed to open image: %w", err)
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return fmt.Errorf("failed to decode image: %w", err)
    }

    resized := resize.Resize(uint(imageJob.Width), uint(imageJob.Height), img, resize.Lanczos3)

    out, err := os.Create(imageJob.OutputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer out.Close()

    err = jpeg.Encode(out, resized, nil)
    if err != nil {
        return fmt.Errorf("failed to encode resized image: %w", err)
    }

    log.Printf("Image resized: %s -> %s", imageJob.InputPath, imageJob.OutputPath)
    return nil
}

func HandleDataProcessingJob(ctx context.Context, job *models.Job) error {
    var dataJob models.DataProcessingJob
    err := json.Unmarshal(job.Payload, &dataJob)
    if err != nil {
        return fmt.Errorf("failed to unmarshal data processing job: %w", err)
    }

    result := processData(dataJob.Data)

    err =storage.StoreResult(result); 
    if err != nil {
        return fmt.Errorf("failed to store result: %w", err)
    }

    log.Printf("Data processed and result stored")
    return nil
}

func processData(data []int) int {
    sum := 0
    for _, v := range data {
        sum += v
    }
    return sum
}
