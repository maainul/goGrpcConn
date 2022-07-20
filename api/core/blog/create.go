package blog

import (
	"context"
	"goGrpcConn/api/storage"
	"goGrpcConn/svcUtils/logging"
)

func (stgBlg *BlgStrgFuncs) CreateBlog(ctx context.Context, blg storage.Blog) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "CreateBlog")
	id, err := stgBlg.st.CreateBlog(ctx, blg)
	if err != nil {
		logging.WithError(err, log).Error("Failed to create blog")
	}
	return id, nil
}
