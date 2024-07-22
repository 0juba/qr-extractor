package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_welcomeHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: `happy path`,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, `http://localhost:8080/`, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			welcomeHandler := createWelcomeHandler(nil, nil, nil)
			welcomeHandler(tt.args.w, tt.args.r)
		})
	}
}
