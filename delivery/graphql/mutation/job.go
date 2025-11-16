package mutation

import (
	"context"
	_dataloader "jobqueue/delivery/graphql/dataloader"
	"jobqueue/delivery/graphql/resolver"
	_interface "jobqueue/interface"
)

type JobMutation struct {
	jobService _interface.JobService
	dataloader *_dataloader.GeneralDataloader
}

// Enqueue resolves the Enqueue mutation
func (q JobMutation) Enqueue(ctx context.Context, args struct { // Change args type to struct for input
	Task string
}) (*resolver.JobResolver, error) {
	// Call the service layer to enqueue the job
	jobID, err := q.jobService.Enqueue(ctx, args.Task)
	if err != nil {
		return nil, err
	}

	// Fetch the created job to return it to the client
	job, err := q.jobService.GetJobByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	return &resolver.JobResolver{
		Data:       *job,
		JobService: q.jobService,
		Dataloader: q.dataloader,
	}, nil
}

// NewJobMutation to create new instance
func NewJobMutation(jobService _interface.JobService, dataloader *_dataloader.GeneralDataloader) JobMutation {
	return JobMutation{
		jobService: jobService,
		dataloader: dataloader,
	}
}
