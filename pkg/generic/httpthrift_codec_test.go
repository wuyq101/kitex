package generic

import (
	"bytes"
	stdjson "encoding/json"
	"net/http"
	"testing"

	"github.com/cloudwego/kitex/internal/test"
)

func TestFromHTTPRequest(t *testing.T) {
	jsonBody := `{"a": 1111111111111, "b": "hello"}`
	req, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewBufferString(jsonBody))
	test.Assert(t, err == nil)
	customReq, err := FromHTTPRequest(req)
	test.Assert(t, err == nil)
	test.DeepEqual(t, customReq.Body, map[string]interface{}{
		"a": stdjson.Number("1111111111111"),
		"b": "hello",
	})
}
