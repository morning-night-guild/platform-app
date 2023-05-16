package external_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/morning-night-guild/platform-app/integration/helper"
)

func TestExternal(t *testing.T) {
	t.Parallel()

	url := helper.GetFirebaseURL(t)

	t.Run("Firebaseに接続できる", func(t *testing.T) {
		t.Parallel()

		res, err := http.DefaultClient.Get(fmt.Sprintf("%s/emulator/openapi.json", url))
		if err != nil {
			t.Fatalf("failed to connect to firebase: %v", err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("failed to connect to firebase: %v", res.Status)
		}
	})
}
