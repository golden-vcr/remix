package queries_test

import (
	"context"
	"testing"

	"github.com/golden-vcr/remix/gen/queries"
	"github.com/golden-vcr/server-common/querytest"
	"github.com/stretchr/testify/assert"
)

func Test_SyncClip(t *testing.T) {
	tx := querytest.PrepareTx(t)
	q := queries.New(tx)

	querytest.AssertCount(t, tx, 0, "SELECT COUNT(*) FROM remix.clip")

	err := q.SyncClip(context.Background(), queries.SyncClipParams{
		ID:       "you_are_the_one",
		Title:    "God's Property: You are the One",
		Duration: 215,
		TapeID:   17,
	})
	assert.NoError(t, err)
	querytest.AssertCount(t, tx, 1, `
		SELECT COUNT(*) FROM remix.clip
			WHERE id = 'you_are_the_one'
			AND title = 'God''s Property: You are the One'
			AND duration = 215
			AND tape_id = 17
	`)

	err = q.SyncClip(context.Background(), queries.SyncClipParams{
		ID:       "you_are_the_one",
		Title:    "Kirk! You it",
		Duration: 216,
		TapeID:   17,
	})
	assert.NoError(t, err)
	querytest.AssertCount(t, tx, 1, `
		SELECT COUNT(*) FROM remix.clip
			WHERE id = 'you_are_the_one'
			AND title = 'Kirk! You it'
			AND duration = 216
			AND tape_id = 17
	`)

	querytest.AssertCount(t, tx, 1, "SELECT COUNT(*) FROM remix.clip")
}
