package repository

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sourcepods/sourcepods/pkg/storage"
)

//LoggingRequestID returns the request ID as string for logging
type LoggingRequestID func(context.Context) string

type loggingService struct {
	service   Service
	requestID LoggingRequestID
	logger    log.Logger
}

// NewLoggingService wraps the Service and provides logging for its methods.
func NewLoggingService(s Service, requestID LoggingRequestID, logger log.Logger) Service {
	return &loggingService{service: s, requestID: requestID, logger: logger}
}

func (s *loggingService) List(ctx context.Context, owner string) ([]*Repository, string, error) {
	start := time.Now()

	repositories, owner, err := s.service.List(ctx, owner)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "List",
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log(
			"msg", "failed to list repositories by owner's username",
			"err", err,
		)
	} else {
		level.Debug(logger).Log()
	}

	return repositories, owner, err

}

func (s *loggingService) Find(ctx context.Context, owner string, name string) (*Repository, string, error) {
	start := time.Now()

	repository, owner, err := s.service.Find(ctx, owner, name)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "Find",
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log(
			"msg", "failed to find repository by owner & name",
			"err", err,
		)
	} else {
		level.Debug(logger).Log()
	}

	return repository, owner, err
}

func (s *loggingService) Create(ctx context.Context, owner string, repository *Repository) (*Repository, error) {
	start := time.Now()

	repository, err := s.service.Create(ctx, owner, repository)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "Create",
		"duration", time.Since(start),
	)

	if err != nil {
		if err != ErrAlreadyExists {
			level.Warn(logger).Log(
				"msg", "failed to create repository",
				"err", err,
			)
		}
	} else {
		level.Debug(logger).Log(
			"owner", owner,
			"name", repository.Name,
		)
	}

	return repository, err
}
func (s *loggingService) Branches(ctx context.Context, owner string, name string) ([]*Branch, error) {
	start := time.Now()

	branches, err := s.service.Branches(ctx, owner, name)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "Branches",
		"duration", time.Since(start),
	)

	if err != nil {
		if err != ErrAlreadyExists {
			level.Warn(logger).Log(
				"msg", "failed to list branches",
				"err", err,
			)
		}
	} else {
		level.Debug(logger).Log(
			"owner", owner,
			"name", name,
		)
	}

	return branches, err
}

func (s *loggingService) Commit(ctx context.Context, owner string, name string, rev string) (storage.Commit, error) {
	start := time.Now()

	commit, err := s.service.Commit(ctx, owner, name, rev)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "Commit",
		"owner", owner,
		"name", name,
		"rev", rev,
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log(
			"msg", "failed to get the commit for repository",
			"err", err,
		)
	} else {
		level.Debug(logger).Log()
	}

	return commit, err
}

func (s *loggingService) ListCommits(ctx context.Context, owner, name, rev string, limit, skip int64) ([]storage.Commit, error) {
	start := time.Now()

	commits, err := s.service.ListCommits(ctx, owner, name, rev, limit, skip)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "Commit",
		"owner", owner,
		"name", name,
		"rev", rev,
		"limit", limit,
		"skip", skip,
		"duration", time.Since(start),
	)

	if err != nil {
		level.Warn(logger).Log(
			"msg", "failed to get the commit for repository",
			"err", err,
		)
	} else {
		level.Debug(logger).Log()
	}

	return commits, err
}

func (s *loggingService) Tree(ctx context.Context, owner, name, rev, path string) ([]storage.TreeEntry, error) {
	start := time.Now()

	tree, err := s.service.Tree(ctx, owner, name, rev, path)

	logger := log.With(s.logger,
		"request", s.requestID(ctx),
		"method", "Tree",
		"owner", owner,
		"name", name,
		"rev", rev,
		"path", path,
		"duration", time.Since(start),
	)

	if err != nil && err != ErrRepositoryNotFound {
		level.Warn(logger).Log(
			"msg", "failed to get the tree for repository",
			"err", err,
		)
	} else {
		level.Debug(logger).Log()
	}

	return tree, err
}
