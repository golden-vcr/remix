package state

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

func Test_Server_handleGetClips(t *testing.T) {
	tests := []struct {
		name       string
		q          *mockQueries
		wantStatus int
		wantBody   string
	}{
		{
			"normal usage",
			&mockQueries{
				clips: []queries.GetClipsRow{
					{
						ID:       "my_cool_clip",
						Title:    "My Very Cool Clip",
						Duration: int32(140),
						TapeID:   int32(42),
					},
				},
			},
			http.StatusOK,
			`{"clips":[{"id":"my_cool_clip","title":"My Very Cool Clip","duration":140,"tapeId":42}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				s := &Server{
					q: tt.q,
				}
				req := httptest.NewRequest(http.MethodGet, "/clips", nil)
				res := httptest.NewRecorder()
				s.handleGetClips(res, req)

				b, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				body := strings.TrimSuffix(string(b), "\n")
				assert.Equal(t, tt.wantStatus, res.Code)
				assert.Equal(t, tt.wantBody, body)
			})
		})
	}
}

type mockQueries struct {
	clips []queries.GetClipsRow
}

func (m *mockQueries) GetClips(ctx context.Context) ([]queries.GetClipsRow, error) {
	return m.clips, nil
}
