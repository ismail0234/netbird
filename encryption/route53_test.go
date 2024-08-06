package encryption

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRoute53TLSConfig(t *testing.T) {
	t.SkipNow() // This test requires AWS credentials
	exampleString := "Hello, world!"
	rtls := &Route53TLS{
		DataDir:            t.TempDir(),
		Email:              os.Getenv("AWS_EMAIL"),
		AwsAccessKeyID:     os.Getenv("AWS_KEY"),
		AwsSecretAccessKey: os.Getenv("AWS_SECRET"),
		Domains:            []string{os.Getenv("DOMAIN")},
	}
	tlsConfig, err := rtls.GetCertificate()
	if err != nil {
		t.Errorf("Route53TLSConfig failed: %v", err)
	}

	server := &http.Server{
		Addr: ":8443",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(exampleString))
		}),
		TLSConfig: tlsConfig,
	}

	go func() {
		err := server.ListenAndServeTLS("", "")
		if err != http.ErrServerClosed {
			t.Errorf("Failed to start server: %v", err)
		}
	}()
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			t.Errorf("Failed to shutdown server: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)
	resp, err := http.Get("https://relay.godevltd.com:8443")
	if err != nil {
		t.Errorf("Failed to get response: %v", err)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	if string(body) != exampleString {
		t.Errorf("Unexpected response: %s", body)
	}
}
