package state

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/golden-vcr/remix"
	"github.com/golden-vcr/remix/gen/queries"
	"github.com/gorilla/mux"
)

type Queries interface {
	GetClips(ctx context.Context) ([]queries.GetClipsRow, error)
}

type Server struct {
	q Queries
}

func NewServer(q *queries.Queries) *Server {
	return &Server{
		q: q,
	}
}

func (s *Server) RegisterRoutes(r *mux.Router) {
	r.Path("/clips").Methods("GET").HandlerFunc(s.handleGetClips)
}

func (s *Server) handleGetClips(res http.ResponseWriter, req *http.Request) {
	// Query the DB to get an ordered listing of all clips
	rows, err := s.q.GetClips(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare an equivalent list of Clip structs
	clips := make([]remix.Clip, 0, len(rows))
	for _, row := range rows {
		clips = append(clips, remix.Clip{
			Id:       row.ID,
			Title:    row.Title,
			Duration: int(row.Duration),
			TapeId:   int(row.TapeID),
		})
	}

	// Return our result object, JSON-serialized
	result := remix.ClipListing{Clips: clips}
	if err := json.NewEncoder(res).Encode(result); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
