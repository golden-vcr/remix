package admin

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golden-vcr/remix/gen/queries"
	"github.com/stretchr/testify/assert"
)

func Test_Server_handleSetTape(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		q             *mockQueries
		wantStatus    int
		wantBody      string
		wantSyncCalls []queries.SyncClipParams
	}{
		{
			"normal usage",
			`{"id":"my_cool_clip","title":"My Very Cool Clip","duration":140,"tapeId":42}`,
			&mockQueries{},
			http.StatusNoContent,
			"",
			[]queries.SyncClipParams{
				{
					ID:       "my_cool_clip",
					Title:    "My Very Cool Clip",
					Duration: int32(140),
					TapeID:   int32(42),
				},
			},
		},
		{
			"invalid tapeId is a 400 error",
			`{"id":"my_cool_clip","title":"My Very Cool Clip","duration":140,"tapeId":0}`,
			&mockQueries{},
			http.StatusBadRequest,
			"'tapeId' value is required",
			nil,
		},
		{
			"invalid duration is a 400 error",
			`{"id":"my_cool_clip","title":"My Very Cool Clip","duration":-99,"tapeId":42}`,
			&mockQueries{},
			http.StatusBadRequest,
			"'duration' value is required",
			nil,
		},
		{
			"missing title is a 400 error",
			`{"id":"my_cool_clip","duration":140,"tapeId":42}`,
			&mockQueries{},
			http.StatusBadRequest,
			"'title' value is required",
			nil,
		},
		{
			"missing id is a 400 error",
			`{"id":"","title":"My Very Cool Clip","duration":140,"tapeId":42}`,
			&mockQueries{},
			http.StatusBadRequest,
			"'id' value is required",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				s := &Server{
					q: tt.q,
				}
				req := httptest.NewRequest(http.MethodPost, "/admin/clip", strings.NewReader(tt.body))
				res := httptest.NewRecorder()
				s.handlePostClip(res, req)

				b, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				body := strings.TrimSuffix(string(b), "\n")
				assert.Equal(t, tt.wantStatus, res.Code)
				assert.Equal(t, tt.wantBody, body)
				assert.Equal(t, tt.wantSyncCalls, tt.q.syncCalls)
			})
		})
	}
}

type mockQueries struct {
	syncCalls []queries.SyncClipParams
}

func (m *mockQueries) SyncClip(ctx context.Context, arg queries.SyncClipParams) error {
	m.syncCalls = append(m.syncCalls, arg)
	return nil
}
