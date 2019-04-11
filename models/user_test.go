package models

import (
	"reflect"
	"testing"
	"time"
)

func TestUser_MarshalJSON(t *testing.T) {
	type fields struct {
		UUID      string
		Email     string
		Password  string
		Name      string
		CreatedAt time.Time
		DeletedAt *time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Test Marshall",
			fields: fields{
				UUID:     "abc",
				Email:    "test@test.com",
				Name:     "Testy",
				Password: "abasjdlsjfsdjflkdjsf",
			},
			want: []byte(`{"uuid":"abc","email":"test@test.com","name":"Testy","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				UUID:      tt.fields.UUID,
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
				Name:      tt.fields.Name,
				CreatedAt: tt.fields.CreatedAt,
				DeletedAt: tt.fields.DeletedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			got, err := u.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
