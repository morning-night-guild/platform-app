package model_test

import (
	"crypto/rsa"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

func TestNewSession(t *testing.T) {
	t.Parallel()

	type args struct {
		sessionID auth.SessionID
		userID    user.UserID
		publicKey rsa.PublicKey
		issuedAt  time.Time
		expiresAt time.Time
	}

	now := time.Now()

	tests := []struct {
		name    string
		args    args
		want    model.Session
		wantErr bool
	}{
		{
			name: "セッションが作成できる",
			args: args{
				sessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				userID:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				publicKey: rsa.PublicKey{},
				issuedAt:  now,
				expiresAt: now.Add(time.Hour * 24 * 30),
			},
			want: model.Session{
				SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				UserID:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				PublicKey: rsa.PublicKey{},
				IssuedAt:  now,
				ExpiresAt: now.Add(time.Hour * 24 * 30),
			},
			wantErr: false,
		},
		{
			name: "セッションが作成できない",
			args: args{
				sessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				userID:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				publicKey: rsa.PublicKey{},
				issuedAt:  now.Add(time.Hour),
				expiresAt: now,
			},
			want:    model.Session{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := model.NewSession(tt.args.sessionID, tt.args.userID, tt.args.publicKey, tt.args.issuedAt, tt.args.expiresAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSession() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionIsExpired(t *testing.T) {
	t.Parallel()

	type fields struct {
		SessionID auth.SessionID
		UserID    user.UserID
		PublicKey rsa.PublicKey
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
				SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				UserID:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				PublicKey: rsa.PublicKey{},
				IssuedAt:  now,
				ExpiresAt: now.Add(-time.Hour * 24 * 30),
			},
			want: true,
		},
		{
			name: "有効期限内",
			fields: fields{
				SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				UserID:    user.UserID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				PublicKey: rsa.PublicKey{},
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
			sss := model.Session{
				SessionID: tt.fields.SessionID,
				UserID:    tt.fields.UserID,
				PublicKey: tt.fields.PublicKey,
				IssuedAt:  tt.fields.IssuedAt,
				ExpiresAt: tt.fields.ExpiresAt,
			}
			if got := sss.IsExpired(); got != tt.want {
				t.Errorf("Session.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
