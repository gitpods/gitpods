package storage

import (
	"context"
	"io"
	"sync"
	"time"

	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"gitlab.com/gitlab-org/gitaly/streamio"
	"google.golang.org/grpc"
)

// Client holds the gRPC-connection to the storage-server
type Client struct {
	repos    RepositoryClient
	branches BranchClient
	commits  CommitClient
	ssh      SSHClient
}

// NewClient returns a new Storage client.
func NewClient(storageAddr string) (*Client, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(grpcopentracing.UnaryClientInterceptor()))
	opts = append(opts, grpc.WithStreamInterceptor(grpcopentracing.StreamClientInterceptor()))

	conn, err := grpc.DialContext(context.Background(), storageAddr, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		repos:    NewRepositoryClient(conn),
		branches: NewBranchClient(conn),
		commits:  NewCommitClient(conn),
		ssh:      NewSSHClient(conn),
	}, nil
}

// Create a repository
func (c *Client) Create(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Create")
	span.SetTag("id", id)
	defer span.Finish()

	_, err := c.repos.Create(ctx, &CreateRequest{Id: id})
	return err
}

// SetDescription of a repository
func (c *Client) SetDescription(ctx context.Context, id, description string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Description")
	span.SetTag("id", id)
	span.SetTag("description", description)
	defer span.Finish()

	_, err := c.repos.SetDescriptions(ctx, &SetDescriptionRequest{
		Id:          id,
		Description: description,
	})
	return err
}

// Branches returns all branches of a repository
func (c *Client) Branches(ctx context.Context, id string) ([]Branch, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Branches")
	span.SetTag("id", id)
	defer span.Finish()

	res, err := c.branches.List(ctx, &BranchesRequest{Id: id})

	var branches []Branch
	for _, b := range res.Branch {
		branches = append(branches, Branch{
			Name: b.Name,
			Sha1: b.Sha1,
			Type: b.Type,
		})
	}

	return branches, err
}

// Commit returns a single commit from a given repository
func (c *Client) Commit(ctx context.Context, id, ref string) (Commit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Commit")
	span.SetTag("id", id)
	span.SetTag("ref", ref)
	defer span.Finish()

	req := &CommitRequest{
		Id:  id,
		Ref: ref,
	}

	res, err := c.commits.Get(ctx, req)
	if err != nil {
		return Commit{}, err
	}

	return Commit{
		Hash:    res.GetHash(),
		Tree:    res.GetTree(),
		Parent:  res.GetParent(),
		Message: res.GetMessage(),
		Author: Signature{
			Name:  res.GetAuthor(),
			Email: res.GetAuthorEmail(),
			Date:  time.Unix(res.GetAuthorDate(), 0),
		},
		Committer: Signature{
			Name:  res.GetCommitter(),
			Email: res.GetCommitterEmail(),
			Date:  time.Unix(res.GetCommitterDate(), 0),
		},
	}, nil
}

//Tree returns the files and folders at a given ref at a path in a repository
func (c *Client) Tree(ctx context.Context, id, ref, path string) ([]TreeEntry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Tree")
	span.SetTag("repo_path", id)
	span.SetTag("ref", ref)
	span.SetTag("path", path)
	defer span.Finish()

	req := &TreeRequest{
		Id:   id,
		Ref:  ref,
		Path: path,
	}

	res, err := c.repos.Tree(ctx, req)
	if err != nil {
		return nil, err
	}

	var treeEntries []TreeEntry
	for _, te := range res.TreeEntries {
		treeEntries = append(treeEntries, TreeEntry{
			Mode:   te.Mode,
			Type:   te.Type,
			Object: te.Object,
			Path:   te.Path,
		})
	}

	return treeEntries, nil
}

// UploadPack to a git-repo
func (c *Client) UploadPack(ctx context.Context, id string, stdin io.Reader, stdout, stderr io.Writer) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.UploadPack")
	span.SetTag("repo_path", id)
	defer span.Finish()

	req := &GRERequest{Id: id}

	stream, err := c.ssh.UploadPack(ctx)
	if err != nil {
		return 0, err
	}

	if err := stream.Send(req); err != nil {
		return 0, nil
	}

	wg := sync.WaitGroup{}

	// Go-routine for sending stdin
	wg.Add(1)
	go func() {
		in := streamio.NewWriter(func(p []byte) error {
			return stream.Send(&GRERequest{Stdin: p})
		})
		// _, err := io.Copy(in, stdin)
		io.Copy(in, stdin)
		stream.CloseSend()
		// Handle err...
		wg.Done()
	}()

	var (
		resp *GREResponse
	)
	for ; err == nil; resp, err = stream.Recv() {
		if resp.GetExitCode() != nil {
			if value := resp.GetExitCode().GetExitCode(); value != 0 {
				return value, nil
			}
			return 0, nil
		}
		if len(resp.GetStderr()) > 0 {
			stderr.Write(resp.GetStderr())
		}
		if len(resp.GetStdout()) > 0 {
			stdout.Write(resp.GetStdout())
		}
	}
	if err == io.EOF {
		err = nil
	}

	wg.Wait()

	return 0, err
}

// ReceivePack from a git-repo
func (c *Client) ReceivePack(ctx context.Context, id string, stdin io.Reader, stdout, stderr io.Writer) (int32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.ReceivePack")
	span.SetTag("repo_path", id)
	defer span.Finish()

	req := &GRERequest{Id: id}

	stream, err := c.ssh.ReceivePack(ctx)
	if err != nil {
		return 0, err
	}

	if err := stream.Send(req); err != nil {
		return 0, nil
	}

	wg := sync.WaitGroup{}

	// Go-routine for sending stdin
	wg.Add(1)
	go func() {
		in := streamio.NewWriter(func(p []byte) error {
			return stream.Send(&GRERequest{Stdin: p})
		})
		// _, err := io.Copy(in, stdin)
		io.Copy(in, stdin)
		stream.CloseSend()
		// Handle err...
		wg.Done()
	}()

	var (
		resp *GREResponse
	)
	for ; err == nil; resp, err = stream.Recv() {
		if resp.GetExitCode() != nil {
			if value := resp.GetExitCode().GetExitCode(); value != 0 {
				return 1, nil
			}
			return 0, nil
		}
		if len(resp.GetStderr()) > 0 {
			if _, err = stderr.Write(resp.GetStderr()); err != nil {
				break
			}
		}
		if len(resp.GetStdout()) > 0 {
			if _, err = stdout.Write(resp.GetStdout()); err != nil {
				break
			}
		}
	}
	if err == io.EOF {
		err = nil
	}

	wg.Wait()

	return 0, err
}
