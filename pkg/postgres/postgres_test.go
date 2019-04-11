package postgres

import (
	"database/sql"
	"reflect"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/kylegrantlucas/platform-exercise/models"
	_ "github.com/lib/pq"
)

func TestCreateDatabase(t *testing.T) {
	type args struct {
		host     string
		port     string
		user     string
		password string
		dbName   string
	}
	tests := []struct {
		name    string
		args    args
		want    *DatabaseConnection
		wantErr bool
	}{
		{
			name: "creating db with empty env var",
			args: args{
				user:     "test",
				password: "",
				host:     "bad",
				port:     "1234",
				dbName:   "test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDatabase(tt.args.host, tt.args.port, tt.args.user, tt.args.password, tt.args.dbName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseConnection_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error building sqlmock: %v", err)
	}

	type fields struct {
		Connection *sql.DB
	}
	type args struct {
		email             string
		name              string
		plaintextPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "valid user create",
			fields: fields{
				Connection: db,
			},
			args: args{
				email:             "test@test.com",
				name:              "testy testerson",
				plaintextPassword: "completelytestpassword",
			},
			want: models.User{
				Email: "test@test.com",
				Name:  "testy testerson",
			},
		},
	}
	for _, tt := range tests {
		mock.ExpectQuery(regexp.QuoteMeta(queries["create_user"])).WillReturnRows(sqlmock.NewRows([]string{"uuid", "email", "name", "created_at", "updated_at"}).AddRow("abc", "test@test.com", "testy testerson", time.Now(), time.Now()))

		t.Run(tt.name, func(t *testing.T) {
			d := &DatabaseConnection{
				Connection: tt.fields.Connection,
			}
			got, err := d.CreateUser(tt.args.email, tt.args.name, tt.args.plaintextPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("DatabaseConnection.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.want.CreatedAt = got.CreatedAt
			tt.want.UpdatedAt = got.UpdatedAt
			tt.want.UUID = got.UUID

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DatabaseConnection.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseConnection_UpdateUserByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error building sqlmock: %v", err)
	}

	currentTime := time.Now()

	type fields struct {
		Connection *sql.DB
	}
	type args struct {
		uuid              string
		email             string
		name              string
		plaintextPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "valid user update",
			fields: fields{
				Connection: db,
			},
			args: args{
				uuid:              "abc",
				email:             "test@test.com",
				name:              "testy testerson",
				plaintextPassword: "completelytestpassword",
			},
			want: models.User{
				UUID:      "abc",
				Email:     "test@test.com",
				Name:      "testy testerson",
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
		},
	}
	for _, tt := range tests {
		mock.ExpectQuery(regexp.QuoteMeta("update users set email=$1,name=$2,password=$3 where uuid=$4 AND deleted_at IS NULL returning uuid, email, name, created_at, updated_at;")).WillReturnRows(sqlmock.NewRows([]string{"uuid", "email", "name", "created_at", "updated_at"}).AddRow("abc", "test@test.com", "testy testerson", currentTime, currentTime))

		t.Run(tt.name, func(t *testing.T) {
			d := &DatabaseConnection{
				Connection: tt.fields.Connection,
			}
			got, err := d.UpdateUserByUUID(tt.args.uuid, tt.args.email, tt.args.name, tt.args.plaintextPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("DatabaseConnection.UpdateUserByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DatabaseConnection.UpdateUserByUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseConnection_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error building sqlmock: %v", err)
	}

	currentTime := time.Now()

	type fields struct {
		Connection *sql.DB
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "valid user create",
			fields: fields{
				Connection: db,
			},
			args: args{
				email: "test@test.com",
			},
			want: models.User{
				UUID:      "abc",
				Email:     "test@test.com",
				Name:      "testy testerson",
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
				Password:  "abc",
			},
		},
	}
	for _, tt := range tests {
		mock.ExpectQuery(regexp.QuoteMeta(queries["get_user_by_email"])).WillReturnRows(sqlmock.NewRows([]string{"uuid", "email", "name", "created_at", "updated_at", "password"}).AddRow("abc", "test@test.com", "testy testerson", currentTime, currentTime, "abc"))

		t.Run(tt.name, func(t *testing.T) {
			d := &DatabaseConnection{
				Connection: tt.fields.Connection,
			}
			got, err := d.GetUserByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("DatabaseConnection.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DatabaseConnection.GetUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseConnection_GetUserByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error building sqlmock: %v", err)
	}

	currentTime := time.Now()

	type fields struct {
		Connection *sql.DB
	}
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "valid user create",
			fields: fields{
				Connection: db,
			},
			args: args{
				uuid: "abc",
			},
			want: models.User{
				UUID:      "abc",
				Email:     "test@test.com",
				Name:      "testy testerson",
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
				Password:  "abc",
			},
		},
	}
	for _, tt := range tests {
		mock.ExpectQuery(regexp.QuoteMeta(queries["get_user_by_uuid"])).WillReturnRows(sqlmock.NewRows([]string{"uuid", "email", "name", "created_at", "updated_at", "password"}).AddRow("abc", "test@test.com", "testy testerson", currentTime, currentTime, "abc"))

		t.Run(tt.name, func(t *testing.T) {
			d := &DatabaseConnection{
				Connection: tt.fields.Connection,
			}
			got, err := d.GetUserByUUID(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DatabaseConnection.GetUserByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DatabaseConnection.GetUserByUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseConnection_SoftDeleteUserByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error building sqlmock: %v", err)
	}

	currentTime := time.Now()

	type fields struct {
		Connection *sql.DB
	}
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "valid user update",
			fields: fields{
				Connection: db,
			},
			args: args{
				uuid: "abc",
			},
			want: models.User{
				UUID:      "abc",
				Email:     "test@test.com",
				Name:      "testy testerson",
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
				DeletedAt: &currentTime,
			},
		},
	}
	for _, tt := range tests {
		mock.ExpectQuery(regexp.QuoteMeta(queries["soft_delete_user_by_uuid"])).WillReturnRows(sqlmock.NewRows([]string{"uuid", "email", "name", "created_at", "updated_at", "deleted_at"}).AddRow("abc", "test@test.com", "testy testerson", currentTime, currentTime, currentTime))
		t.Run(tt.name, func(t *testing.T) {
			d := &DatabaseConnection{
				Connection: tt.fields.Connection,
			}
			got, err := d.SoftDeleteUserByUUID(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DatabaseConnection.SoftDeleteUserByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DatabaseConnection.SoftDeleteUserByUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseConnection_CreateSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error building sqlmock: %v", err)
	}

	currentTime := time.Now()

	type fields struct {
		Connection *sql.DB
	}
	type args struct {
		userUUID  string
		expiresAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Session
		wantErr bool
	}{
		{
			name: "valid user update",
			fields: fields{
				Connection: db,
			},
			args: args{
				userUUID:  "abc",
				expiresAt: currentTime,
			},
			want: models.Session{
				UUID:      "abc",
				UserUUID:  "abc",
				CreatedAt: currentTime,
				ExpiresAt: currentTime,
			},
		},
	}
	for _, tt := range tests {
		mock.ExpectQuery(regexp.QuoteMeta(queries["create_session"])).WillReturnRows(sqlmock.NewRows([]string{"uuid", "user_uuid", "created_at", "expires_at"}).AddRow("abc", "abc", currentTime, currentTime))
		t.Run(tt.name, func(t *testing.T) {
			d := &DatabaseConnection{
				Connection: tt.fields.Connection,
			}
			got, err := d.CreateSession(tt.args.userUUID, tt.args.expiresAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("DatabaseConnection.CreateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DatabaseConnection.CreateSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseConnection_GetSessionByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error building sqlmock: %v", err)
	}

	currentTime := time.Now()

	type fields struct {
		Connection *sql.DB
	}
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Session
		wantErr bool
	}{
		{
			name: "valid user update",
			fields: fields{
				Connection: db,
			},
			args: args{
				uuid: "abc",
			},
			want: models.Session{
				UUID:      "abc",
				UserUUID:  "abc",
				CreatedAt: currentTime,
				ExpiresAt: currentTime,
				DeletedAt: &currentTime,
			},
		},
	}
	for _, tt := range tests {
		mock.ExpectQuery(regexp.QuoteMeta(queries["get_session_by_uuid"])).WillReturnRows(sqlmock.NewRows([]string{"uuid", "user_uuid", "created_at", "expires_at", "deleted_at"}).AddRow("abc", "abc", currentTime, currentTime, currentTime))
		t.Run(tt.name, func(t *testing.T) {
			d := &DatabaseConnection{
				Connection: tt.fields.Connection,
			}
			got, err := d.GetSessionByUUID(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DatabaseConnection.GetSessionByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DatabaseConnection.GetSessionByUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseConnection_SoftDeleteSessionByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error building sqlmock: %v", err)
	}

	type fields struct {
		Connection *sql.DB
	}
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "valid session delete",
			fields: fields{
				Connection: db,
			},
			args: args{
				uuid: "abc",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		mock.ExpectExec(regexp.QuoteMeta(queries["soft_delete_session_by_uuid"])).WillReturnResult(sqlmock.NewResult(1, 1))
		t.Run(tt.name, func(t *testing.T) {
			d := &DatabaseConnection{
				Connection: tt.fields.Connection,
			}
			got, err := d.SoftDeleteSessionByUUID(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DatabaseConnection.SoftDeleteSessionByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DatabaseConnection.SoftDeleteSessionByUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
