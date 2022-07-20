package handler

import (
	"net/http"

	"goGrpcConn/svcUtils/logging"
)

func (s *Server) indexHander(w http.ResponseWriter, r *http.Request) {
	logging.FromContext(r.Context()).WithField("method", "index Hander")
}
