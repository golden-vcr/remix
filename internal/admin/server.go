package admin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/golden-vcr/auth"
	"github.com/golden-vcr/remix"
	"github.com/golden-vcr/remix/gen/queries"
	"github.com/golden-vcr/server-common/entry"
	"github.com/gorilla/mux"
)

type Queries interface {
	SyncClip(ctx context.Context, arg queries.SyncClipParams) error
}

type Server struct {
	q Queries
}

func NewServer(q *queries.Queries) *Server {
	return &Server{
		q: q,
	}
}

func (s *Server) RegisterRoutes(c auth.Client, r *mux.Router) {
	// Require broadcaster access for all admin routes
	r.Use(func(next http.Handler) http.Handler {
		return auth.RequireAccess(c, auth.RoleBroadcaster, next)
	})

	// POST /clip allows the broadcaster to sync the details of a clip: this'll create a
	// new clip if the ID is brand new; or modify an existing clip's details if a clip
	// already exists with the given ID
	r.Path("/clip").Methods("POST").HandlerFunc(s.handlePostClip)
}

func (s *Server) handlePostClip(res http.ResponseWriter, req *http.Request) {
	// Parse and validate the request clip
	var clip remix.Clip
	if err := json.NewDecoder(req.Body).Decode(&clip); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if clip.Id == "" {
		http.Error(res, "'id' value is required", http.StatusBadRequest)
		return
	}
	if clip.Title == "" {
		http.Error(res, "'title' value is required", http.StatusBadRequest)
		return
	}
	if clip.Duration <= 0 {
		http.Error(res, "'duration' value is required", http.StatusBadRequest)
		return
	}
	if clip.TapeId <= 0 {
		http.Error(res, "'tapeId' value is required", http.StatusBadRequest)
		return
	}

	// Update the DB to reflect our new desired state for this clip
	if err := s.q.SyncClip(req.Context(), queries.SyncClipParams{
		ID:       clip.Id,
		Title:    clip.Title,
		Duration: int32(clip.Duration),
		TapeID:   int32(clip.TapeId),
	}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	entry.Log(req).Info("Synced clip details", "clip", clip)
	res.WriteHeader(http.StatusNoContent)
}
