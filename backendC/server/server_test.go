package server

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// runTestQuery runs the "query" and unmarshals the JSON body into "output".
//
// If the response body cannot be unmarshaled the test fails.
func runTestQuery(t *testing.T, query string, output interface{}) {
	req := httptest.NewRequest("GET", "/?q="+query, nil)
	res := httptest.NewRecorder()

	Handler().ServeHTTP(res, req)

	err := json.Unmarshal(res.Body.Bytes(), output)
	if err != nil {
		t.Fatalf("unable to unmarshal response into %v: %v", output, err)
	}
}
