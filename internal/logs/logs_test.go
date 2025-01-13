package logs

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http/httptest"
	"testing"
	"testing/slogtest"
)

func TestLol(t *testing.T) {
	// todo nuke or change name
	c := ContextKeyRequestMethod()
	c.LogKey = "test"
	assert.NotEqual(t, "test", contextKeyRequestMethod.LogKey)
}

func TestContextKey_Get_Returns_Error_If_ContextKey_Not_Found(t *testing.T) {
	_, err := contextKeyRequestMethod.Get(context.Background())
	assert.Error(t, err)
}

func TestContextKey_Get_Returns_Value_If_ContextKey_Found(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextKeyRequestMethod, "value")

	v, err := contextKeyRequestMethod.Get(ctx)
	assert.NoError(t, err)
	assert.Equal(t, v, "value")
}

func TestContextKey_AddToContext_Adds_Values_To_Context(t *testing.T) {
	ctx := context.Background()
	ctx = contextKeyRequestMethod.AddToContext(ctx, "GET")
	ctx = contextKeyRequestPath.AddToContext(ctx, "/")

	assert.Equal(t, "GET", ctx.Value(contextKeyRequestMethod).(string))
	assert.Equal(t, "/", ctx.Value(contextKeyRequestPath).(string))
}

func TestContextKey_AddToRecord_Adds_Attrs_To_Record(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req = AddRequestAttributes(req)
	ctx := req.Context()

	record := slog.Record{}
	record = contextKeyRequestMethod.AddToRecord(ctx, record)
	record = contextKeyRequestPath.AddToRecord(ctx, record)

	assert.Equal(t, 2, record.NumAttrs())
}

func Test_AddRequestAttributes_Adds_Request_Attrs_To_Context(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req = AddRequestAttributes(req)

	v, err := contextKeyRequestPath.Get(req.Context())
	assert.NoError(t, err)
	assert.Equal(t, "/", v)

	v, err = contextKeyRequestMethod.Get(req.Context())
	assert.NoError(t, err)
	assert.Equal(t, "GET", v)
}

// See https://github.com/golang/example/blob/master/slog-handler-guide/README.md#testing
func TestContextualHandler_Follows_Contract(t *testing.T) {
	var buf bytes.Buffer
	h := NewContextualHandler(slog.NewJSONHandler(&buf, nil))
	err := slogtest.TestHandler(h, func() []map[string]any {
		return parseJSONLogEntries(t, buf.Bytes())
	})
	if err != nil {
		t.Error(err)
	}
}

func TestContextualHandler_Handle_Adds_Context_To_Record(t *testing.T) {
	var buf bytes.Buffer
	h := NewContextualHandler(slog.NewJSONHandler(&buf, nil))

	ctx := context.Background()
	ctx = contextKeyRequestMethod.AddToContext(ctx, "GET")
	ctx = contextKeyRequestPath.AddToContext(ctx, "/")
	ctx = contextKeySessionID.AddToContext(ctx, "session-id")

	// add tests for additional context keys

	err := h.Handle(ctx, slog.Record{})
	assert.NoError(t, err)

	entries := parseJSONLogEntries(t, buf.Bytes())
	e := entries[0]
	assert.NotEmpty(t, e)
	assert.Contains(t, e, "method")
	assert.Equal(t, "GET", e["method"])
	assert.Contains(t, e, "path")
	assert.Equal(t, "/", e["path"])
	assert.Equal(t, "session-id", e["session_id"])
}

func parseJSONLogEntries(t *testing.T, data []byte) []map[string]any {
	entries := bytes.Split(data, []byte("\n"))
	entries = entries[:len(entries)-1] // last one is empty
	var ms []map[string]any
	for _, e := range entries {
		var m map[string]any
		if err := json.Unmarshal(e, &m); err != nil {
			t.Fatal(err)
		}
		ms = append(ms, m)
	}
	return ms
}
