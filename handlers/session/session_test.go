package session

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kylegrantlucas/platform-exercise/pkg/postgres"
)

func TestCreate(t *testing.T) {
	postgres.DB = &postgres.DBMock{}

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test success",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/sessions", bytes.NewReader([]byte(`{"email": "test@gmail.com", "password": "test"}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Create(tt.args.w, tt.args.r)
		})
	}
}

func TestDelete(t *testing.T) {
	postgres.DB = &postgres.DBMock{}

	r := httptest.NewRequest("DELETE", "/sessions", nil)
	r.Header.Add("X-Verified-Session-Uuid", "abc")

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test success",
			args: args{
				w: httptest.NewRecorder(),
				r: r,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Delete(tt.args.w, tt.args.r)
		})
	}
}
