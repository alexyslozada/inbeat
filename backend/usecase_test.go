package backend

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUseCase_cleanUserName(t *testing.T) {
	table := []struct {
		name     string
		userName string
		want     string
	}{
		{name: "With @", userName: "@alexys_lozada", want: "alexys_lozada"},
		{name: "Without @", userName: "alexys_lozada", want: "alexys_lozada"},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanUserName(tt.userName)
			if got != tt.want {
				t.Fatalf("Got %s, Want %s", got, tt.want)
			}
		})
	}
}

func TestUseCase_doRequest(t *testing.T) {
	table := []struct {
		name       string
		statusCode int
		wantErr    error
		want       string
	}{
		{name: "normal", statusCode: http.StatusOK, wantErr: nil, want: "ok"},
		{name: "normal", statusCode: http.StatusNotFound, wantErr: fmt.Errorf("we get an unexpected status code: %d", http.StatusNotFound), want: ""},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.want))
			}))
			defer svr.Close()
			data, err := doRequest(svr.URL)
			if err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Fatalf("we want %v, got %v", tt.wantErr, err)
				}
			}
			if string(data) != tt.want {
				t.Fatalf("we want body %s, got %s", tt.want, string(data))
			}
		})
	}
}
