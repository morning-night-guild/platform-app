package model_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

func TestNewCode(t *testing.T) {
	t.Parallel()

	type args struct {
		codeID    auth.CodeID
		sessionID auth.SessionID
		issuedAt  time.Time
		expiresAt time.Time
	}

	now := time.Now()

	tests := []struct {
		name    string
		args    args
		want    model.Code
		wantErr bool
	}{
		{
			name: "コードを作成できる",
			args: args{
				codeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				sessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				issuedAt:  now,
				expiresAt: now.Add(time.Hour),
			},
			want: model.Code{
				CodeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				IssuedAt:  now,
				ExpiresAt: now.Add(time.Hour),
			},
			wantErr: false,
		},
		{
			name: "コードを作成できない",
			args: args{
				codeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				sessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				issuedAt:  now.Add(time.Hour),
				expiresAt: now,
			},
			want:    model.Code{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := model.NewCode(tt.args.codeID, tt.args.sessionID, tt.args.issuedAt, tt.args.expiresAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodeIsExpired(t *testing.T) {
	t.Parallel()

	type fields struct {
		CodeID    auth.CodeID
		SessionID auth.SessionID
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
				CodeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				IssuedAt:  now,
				ExpiresAt: now.Add(-time.Hour),
			},
			want: true,
		},
		{
			name: "有効期限内",
			fields: fields{
				CodeID:    auth.CodeID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				SessionID: auth.SessionID(uuid.MustParse("01234567-0123-0123-0123-0123456789ab")),
				IssuedAt:  now,
				ExpiresAt: now.Add(time.Hour),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cd := model.Code{
				CodeID:    tt.fields.CodeID,
				SessionID: tt.fields.SessionID,
				IssuedAt:  tt.fields.IssuedAt,
				ExpiresAt: tt.fields.ExpiresAt,
			}
			if got := cd.IsExpired(); got != tt.want {
				t.Errorf("Code.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
