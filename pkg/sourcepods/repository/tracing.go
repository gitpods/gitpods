package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/sourcepods/sourcepods/pkg/storage"
)

//TracingRequestID returns the request ID as string for tracing
type TracingRequestID func(context.Context) string

type tracingService struct {
	service   Service
	requestID TracingRequestID
}

// NewTracingService wraps the Service and provides tracing for its methods.
func NewTracingService(s Service, requestID TracingRequestID) Service {
	return &tracingService{service: s, requestID: requestID}
}

func (s *tracingService) List(ctx context.Context, owner string) ([]*Repository, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.List")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	defer span.Finish()

	return s.service.List(ctx, owner)
}

func (s *tracingService) Find(ctx context.Context, owner string, name string) (*Repository, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Find")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	return s.service.Find(ctx, owner, name)
}

func (s *tracingService) Create(ctx context.Context, owner string, repository *Repository) (*Repository, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Create")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	span.SetTag("name", repository.Name)
	defer span.Finish()

	return s.service.Create(ctx, owner, repository)
}

func (s *tracingService) Branches(ctx context.Context, owner string, name string) ([]*Branch, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Branches")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	return s.service.Branches(ctx, owner, name)
}

func (s *tracingService) Commit(ctx context.Context, owner string, name string, rev string) (storage.Commit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Commit")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("rev", rev)
	defer span.Finish()

	return s.service.Commit(ctx, owner, name, rev)
}

func (s *tracingService) ListCommits(ctx context.Context, owner, name, ref string, limit, skip int64) ([]storage.Commit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.ListCommits")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("ref", ref)
	span.SetTag("limit", limit)
	span.SetTag("skip", skip)
	defer span.Finish()

	return s.service.ListCommits(ctx, owner, name, ref, limit, skip)
}

func (s *tracingService) Tree(ctx context.Context, owner, name, rev, path string) ([]storage.TreeEntry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Service.Tree")
	span.SetTag("request", s.requestID(ctx))
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("rev", rev)
	span.SetTag("path", path)
	defer span.Finish()

	return s.service.Tree(ctx, owner, name, rev, path)
}
