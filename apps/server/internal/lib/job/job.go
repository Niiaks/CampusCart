package job

import (
	"github.com/Niiaks/campusCart/internal/config"
	lib "github.com/Niiaks/campusCart/internal/lib/email"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type JobService struct {
	client *asynq.Client
	server *asynq.Server
	logger *zerolog.Logger
}

func NewJobService(cfg *config.Config, logger *zerolog.Logger) *JobService {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.Redis.Address,
	})

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.Address},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6, //important
				"default":  3, //most
				"low":      1, //non urgent
			},
		},
	)
	return &JobService{
		client: client,
		server: server,
		logger: logger,
	}
}

func (j *JobService) Start(emailClient *lib.Client) error {
	mux := asynq.NewServeMux()

	// Register email task handlers
	mux.HandleFunc(TaskEmailWelcome, j.HandleEmailWelcome(emailClient))
	mux.HandleFunc(TaskEmailVerification, j.HandleEmailVerification(emailClient))

	j.logger.Info().Msg("Starting background job server")
	if err := j.server.Start(mux); err != nil {
		return err
	}

	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("Stopping background job server")
	j.server.Shutdown()
	j.client.Close()
}

// Enqueue dispatches a task to the background job queue.
func (j *JobService) Enqueue(task *asynq.Task, opts ...asynq.Option) error {
	info, err := j.client.Enqueue(task, opts...)
	if err != nil {
		j.logger.Error().Err(err).Str("task", task.Type()).Msg("failed to enqueue task")
		return err
	}
	j.logger.Info().
		Str("task", task.Type()).
		Str("id", info.ID).
		Str("queue", info.Queue).
		Msg("task enqueued")
	return nil
}
