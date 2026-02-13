package job

import (
	"github.com/Niiaks/campusCart/internal/config"
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

func (j *JobService) Start() error {
	mux := asynq.NewServeMux()

	//register task handlers here

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
