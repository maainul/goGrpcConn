package blog

import (
	"context"
	gk "goGrpcConn/api/gunk/v1/admin/blog"
	"goGrpcConn/api/storage"
	"goGrpcConn/svcUtils/logging"
	"time"
)

func (s *BlogCoreFuncs) CreateBlog(ctx context.Context, req *gk.CreateBlogRequest) (*gk.CreateBlogResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.Blog.CreateBlog")
	id, err := s.bc.CreateBlog(ctx, storage.Blog{
		Name:      req.Blog.Title,
		CreatedAt: time.Now(),
		CreatedBy: req.Blog.CreatedBy,
	})
	if err != nil {
		logging.WithError(err, log).Error("Failed to create blog")
		return nil, err
	}
	return &gk.CreateBlogResponse{
		ID: id,
	}, nil
}
