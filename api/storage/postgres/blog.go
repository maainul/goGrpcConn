package postgres

import (
	"context"
	"goGrpcConn/api/storage"
	"goGrpcConn/svcUtils/logging"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const crtBlg = `
INSERT INTO blogs (
	name,
	created_at,
	created_by
) VALUES (
	:name,
	:created_at,
	:created_by
) RETURNING
	id
`

func (s *Storage) CreateBlog(ctx context.Context, blg storage.Blog) (string, error) {
	log := logging.FromContext(ctx)
	stmt, err := s.db.PrepareNamed(crtBlg)
	if err != nil {
		logging.WithError(err, log).Error("Failed to create blog")
		return "", err
	}
	defer stmt.Close()
	var id string
	if err := stmt.Get(&id, blg); err != nil {
		logging.WithError(err, log).Error("Failed to insert blog")
		return "", status.Errorf(codes.Internal, "%d", err)

	}
	return id, nil

}
