package model_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

func TestNewAuth(t *testing.T) {
	t.Parallel()

	type args struct {
		authID    user.ID
		userID    user.ID
		issuedAt  time.Time
		expiresAt time.Time
	}

	now := time.Now()

	tests := []struct {
		name    string
		args    args
		want    model.Auth
		wantErr bool
	}{
		{
			name: "認証情報が作成できる",
			args: args{
				authID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				userID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				issuedAt:  now,
				expiresAt: now.Add(time.Hour * 24 * 30),
			},
			want: model.Auth{
				AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				UserID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				IssuedAt:  now,
				ExpiresAt: now.Add(time.Hour * 24 * 30),
			},
			wantErr: false,
		},
		{
			name: "認証情報が作成できない",
			args: args{
				authID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				userID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				issuedAt:  now.Add(time.Hour),
				expiresAt: now,
			},
			want:    model.Auth{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := model.NewAuth(tt.args.authID, tt.args.userID, tt.args.issuedAt, tt.args.expiresAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAuth() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthIsExpired(t *testing.T) {
	t.Parallel()

	type fields struct {
		AuthID    user.ID
		ID        user.ID
		IssuedAt  time.Time
		ExpiresAt time.Time
	}

	now := time.Now()

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "有効期限切れ",
			fields: fields{
				AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				ID:        user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				IssuedAt:  now,
				ExpiresAt: now.Add(-time.Hour * 24 * 30),
			},
			want: true,
		},
		{
			name: "有効期限内",
			fields: fields{
				AuthID:    user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				ID:        user.ID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				IssuedAt:  now,
				ExpiresAt: now.Add(time.Hour * 24 * 30),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			at := model.Auth{
				AuthID:    tt.fields.AuthID,
				UserID:    tt.fields.ID,
				IssuedAt:  tt.fields.IssuedAt,
				ExpiresAt: tt.fields.ExpiresAt,
			}
			if got := at.IsExpired(); got != tt.want {
				t.Errorf("Auth.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
