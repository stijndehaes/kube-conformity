package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	log "github.com/sirupsen/logrus"
)

func TestConstructConfig(t *testing.T) {
	configLocation = "config.yaml"
	config, err := ConstructConfig()
	assert.Nil(t, err)
	assert.Len(t, config.PodRulesLabelsFilledIn, 1)
	assert.Len(t, config.PodRulesLimitsFilledIn, 1)
	assert.Len(t, config.PodRulesRequestsFilledIn, 1)
}

func TestConfigureLogging(t *testing.T) {
	debug = false
	jsonLogging = false

	ConfigureLogging()
	assert.Equal(t, &log.TextFormatter{},log.StandardLogger().Formatter)
	assert.Equal(t, log.InfoLevel,log.StandardLogger().Level)

	debug = true
	jsonLogging = true

	ConfigureLogging()
	assert.Equal(t, &log.JSONFormatter{},log.StandardLogger().Formatter)
	assert.Equal(t, log.DebugLevel,log.StandardLogger().Level)
}

func TestConstructConfig_InvalidLocation(t *testing.T) {
	configLocation = "invalid.yaml"
	_, err := ConstructConfig()
	assert.NotNil(t, err)
}

func Test_configurePrometheus(t *testing.T) {
	config, _ := ConstructConfig()
	configurePrometheus(config)
}

func Test_defaultPageHandler(t *testing.T) {
	config, _ := ConstructConfig()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(defaultPageHandler(config))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}


func Test_healthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
