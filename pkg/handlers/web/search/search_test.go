package search

import (
	"context"
	"homeserver/db"
	"net/url"
	"testing"
)

type pagedFiles struct {
	db.FileRepo
	t *testing.T
}

func (p pagedFiles) FuzzySearch(_ context.Context, query, directory string, owner *db.User, limit int, offset int64) ([]db.UserFile, error) {
	p.t.Helper()
	if query != "a & b" || directory != "alice/my files" || owner.ID != 7 || limit != 21 || offset != 20 {
		p.t.Fatalf("unexpected search: %q %q owner=%d limit=%d offset=%d", query, directory, owner.ID, limit, offset)
	}
	return make([]db.UserFile, 21), nil
}

func TestFileSearchPageKeepsScopeAndQuery(t *testing.T) {
	repos := &db.Database{Files: pagedFiles{t: t}}
	data, err := FileSearchHandler(context.Background(), repos, "a & b", []string{"alice", "my files"}, &db.User{ID: 7}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(data.Entries); got != 20 {
		t.Fatalf("got %d results, want 20", got)
	}
	for link, page := range map[string]string{data.PreviousPageURL: "1", data.NextPageURL: "3"} {
		u, err := url.Parse(link)
		if err != nil {
			t.Fatal(err)
		}
		if u.Path != "/search/my+files" || u.Query().Get("query") != "a & b" || u.Query().Get("page") != page {
			t.Errorf("bad pagination URL: %q", link)
		}
	}
}
