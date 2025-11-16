package query

import (
	"context"
	_dataloader "jobqueue/delivery/graphql/dataloader"
	"jobqueue/delivery/graphql/resolver"
	_interface "jobqueue/interface"
)

type JobQuery struct {
	jobService _interface.JobService
	dataloader *_dataloader.GeneralDataloader
}

// Jobs resolves the Jobs query
func (q JobQuery) Jobs(ctx context.Context) ([]resolver.JobResolver, error) {
	jobs, err := q.jobService.GetAllJobs(ctx)
	if err != nil {
		return nil, err
	}

	resolvers := make([]resolver.JobResolver, len(jobs))
	for i, job := range jobs {
		resolvers[i] = resolver.JobResolver{
			Data:       *job,
			JobService: q.jobService,
			Dataloader: q.dataloader,
		}
	}
	return resolvers, nil
}

// Job resolves the Job(id: String!) query
func (q JobQuery) Job(ctx context.Context, args struct {
	ID string
}) (*resolver.JobResolver, error) {
	job, err := q.jobService.GetJobByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}

	return &resolver.JobResolver{
		Data:       *job,
		JobService: q.jobService,
		Dataloader: q.dataloader,
	}, nil
}

// JobStatus resolves the JobStatus query
func (q JobQuery) JobStatus(ctx context.Context) (resolver.JobStatusResolver, error) {
	jobStatus, err := q.jobService.GetJobStatus(ctx)
	if err != nil {
		return resolver.JobStatusResolver{}, err
	}

	return resolver.JobStatusResolver{
		Data:       jobStatus,
		JobService: q.jobService,
		Dataloader: q.dataloader,
	}, nil
}

func NewJobQuery(jobService _interface.JobService,
	dataloader *_dataloader.GeneralDataloader) JobQuery {
	return JobQuery{
		jobService: jobService,
		dataloader: dataloader,
	}
}
