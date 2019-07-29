package backendtest

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

type Runner struct {
	FixturePath string
	Update      bool
	Handler     http.Handler
}

func (r *Runner) RunQuerySet(t *testing.T, queries []string) {
	for _, q := range queries {
		r.RunQuery(t, q)
	}
}

func (r *Runner) RunQuery(t *testing.T, query string) {
	w := httptest.NewRecorder()

	params := make(url.Values, 1)
	params.Set("q", query)
	path := "/?" + params.Encode()
	req := httptest.NewRequest("GET", path, nil)

	r.Handler.ServeHTTP(w, req)

	actual := w.Body.Bytes()
	if r.Update {
		r.storeQueryFixture(t, query, actual)
	}
	expected := r.loadQueryFixture(t, query)

	if !bytes.Equal(actual, expected) {
		t.Errorf("mismatch for %s", query)

		var expectedBuffer, actualBuffer bytes.Buffer

		json.Indent(&expectedBuffer, expected, "", "  ")
		json.Indent(&actualBuffer, actual, "", "  ")

		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(expectedBuffer.String()),
			B:        difflib.SplitLines(actualBuffer.String()),
			ToFile:   "actual",
			FromFile: "expected",
			Context:  5,
		})
		if err != nil {
			t.Errorf("diff failed: %v", err)
		} else {
			t.Log("\n" + diff)
		}
	}
}

func (r *Runner) storeQueryFixture(t *testing.T, query string, data []byte) {
	path := filepath.Join(r.FixturePath, hashQuery(query))

	err := ioutil.WriteFile(path, data, 0777)
	if err != nil {
		t.Fatalf("error storing fixture: %v", err)
	}
}

func (r *Runner) loadQueryFixture(t *testing.T, query string) []byte {
	path := filepath.Join(r.FixturePath, hashQuery(query))

	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("error loading fixture: %v", err)
	}

	return data
}

func hashQuery(q string) string {
	hash := crc32.NewIEEE()
	hash.Write([]byte(q))

	key := make([]byte, 0, hash.Size())
	key = hash.Sum(key)

	return hex.EncodeToString(key)
}
