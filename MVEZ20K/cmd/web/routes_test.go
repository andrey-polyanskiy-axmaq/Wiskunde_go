package main

import (
	"MVEZ20K/internal/config"
	"fmt"
	"testing"
)

func TestRouters(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing; test passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, type is %T", v))
	}

}
