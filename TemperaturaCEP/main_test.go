package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInvalidCEP(t *testing.T) {
	req := httptest.NewRequest("GET", "/temperatura?cep=123", nil)
	w := httptest.NewRecorder()

	temperaturaHandler(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Esperado 422, obteve %d", w.Code)
	}
}

func TestNotFoundCEP(t *testing.T) {
	req := httptest.NewRequest("GET", "/temperatura?cep=99999999", nil)
	w := httptest.NewRecorder()

	temperaturaHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Esperado 404, obteve %d", w.Code)
	}
}
