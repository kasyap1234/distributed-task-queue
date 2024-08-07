# Distributed Task Queue

A robust, scalable distributed task queue system implemented in Go, designed to handle various types of jobs efficiently.

## Features

- Concurrent job processing
- Redis-based queue for job storage and retrieval
- Job scheduling for recurring tasks
- RESTful API for job submission
- Configurable worker count and timeout
- Exponential backoff strategy for efficient queue polling
- Support for multiple job types (email, image resizing, data processing)
- Persistent storage for job results

## Components

### Worker (worker.go)

The worker component is responsible for processing jobs from the queue. It supports:

- Concurrent job processing
- Custom job type handlers
- Exponential backoff for efficient queue polling

### Scheduler (scheduler.go)

The scheduler allows for the creation of recurring jobs. It uses a cron-like syntax for job scheduling.

### Queue (queue.go)

Implements a Redis-based queue for job storage and retrieval. Supports operations like:

- Enqueue
- Dequeue
- Update job status

### Job (job.go)

Defines the structure for different job types, including:

- Email jobs
- Image resizing jobs
- Data processing jobs

### Handlers (handler.go)

Contains the logic for processing different job types:

- Email sending
- Image resizing
- Data processing

### Storage (storage.go)

Manages persistent storage of job results using PostgreSQL.

### Config (config.go)

Handles application configuration using Viper, allowing for easy management of:

- Redis connection details
- Worker count and timeout
- API server address

## Getting Started

1. Ensure Redis and PostgreSQL are installed and running.
2. Set up the configuration in `config.go` or use environment variables.
3. Run the main application:


4. The API server will start, and you can submit jobs via POST requests to `/submit`.

## API Usage

Submit a job:

POST /submit Content-Type: application/json

{ "type": "email", "payload": { "to": "user@example.com", "subject": "Test Email", "body": "This is a test email." } }


## Configuration

The application can be configured using environment variables:

- `REDISADDR`: Redis server address (default: "127.0.0.1:6379")
- `MAXRETRIES`: Maximum job retry attempts (default: 5)
- `WORKERCOUNT`: Number of concurrent workers (default: 10)
- `APIADDR`: API server address (default: "127.0.0.1:8080")
- `WORKERTIMEOUT`: Worker timeout in seconds (default: 10)

## Extending the System

To add new job types:

1. Define the job structure in `job.go`
2. Implement the job handler in `handler.go`
3. Register the new handler in `main.go`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
