//go:build e2e

package e2e_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
	log "github.com/sirupsen/logrus"
)

const (
	defaultAttempts = 20
	host            = "app:8080"
	healthPath      = "http://" + host + "/health"
	basePath        = "http://" + host
)

func TestMain(m *testing.M) {
	if err := HealthCheck(defaultAttempts); err != nil {
		log.Fatalf("health check failed: %v", err)
	}
	log.Infof("e2e tests: host %s is available", basePath)
	os.Exit(m.Run())
}

func HealthCheck(attempts int) error {
	var err error
	for attempts > 0 {
		err = Do(
			Get(healthPath),
			Expect().Status().Equal(http.StatusOK),
		)

		if err == nil {
			return nil
		}
		time.Sleep(time.Second)
		attempts--
	}
	return err
}
