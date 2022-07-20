package blog

import (
	"context"
	"log"
	"os"
	"path/filepath"
	bcr "goGrpcConn/api/core/blog"
	gk "goGrpcConn/api/gunk/v1/admin/blog"
	"goGrpcConn/api/storage/postgres"
	"testing"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestBlogCoreConn(t *testing.T) {
	dbconn, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)
	svc := BlogCoreConn(bcr.ConnWithStorage(dbconn))
	tests := []struct {
		name    string
		in      gk.Blog
		want    *gk.Blog
		WantErr bool
	}{
		{
			name: "CREATE_BLOG_SUCCESS",
			in: gk.Blog{
				Title:     "New Blog Title",
				CreatedAt: timestamppb.Now(),
				CreatedBy: "12345",
			},
			want: &gk.Blog{
				Title:     "New Blog Title",
				CreatedAt: timestamppb.Now(),
				CreatedBy: "12345",
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			id, err := svc.CreateBlog(context.TODO(), &gk.CreateBlogRequest{
				Blog: &tt.in,
			})
			if err != nil {
				t.Errorf("Storage.CreateBlog() error = %v, wantErr %v", err, tt.WantErr)
				return
			}
			if id.ID == "" {
				log.Fatalf("Create Blog Failed")
			}

		})
	}
}
