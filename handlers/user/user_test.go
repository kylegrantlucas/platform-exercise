package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
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
				r: httptest.NewRequest("POST", "/users", nil),
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

	r := httptest.NewRequest("POST", "/sessions", nil)
	r.Header.Add("X-Verified-User-Uuid", "abc")
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

func TestUpdate(t *testing.T) {
	postgres.DB = &postgres.DBMock{}

	r := httptest.NewRequest("POST", "/sessions", bytes.NewReader([]byte(`{"email": "test", "password": "9X&5eQ#TI9IzBM", "name": "Testers"}`)))
	r.Header.Add("X-Verified-User-Uuid", "abc")
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
			Update(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseUserRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    userRequest
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseUserRequest(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUserRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseUserRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
