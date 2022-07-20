package blog

import (
	"context"
	"goGrpcConn/api/storage"
	gk "goGrpcConn/api/gunk/v1/admin/blog"
)

type BlogCoreFuncs struct {
	bc BlogCore
	gk.UnimplementedBlogServiceServer
}

type BlogCore interface {
	CreateBlog(context.Context, storage.Blog) (string, error)
}

func BlogCoreConn(bc BlogCore) *BlogCoreFuncs {
	return &BlogCoreFuncs{bc: bc}
}
