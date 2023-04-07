package internal

import (
	"crypto/rand"
	"net/http"
	"time"

	"github.com/leg100/surl"
	"github.com/rs/zerolog"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func signURL(url string) (string, error) {
	signer := surl.New([]byte(appConfig.SecretKey))
	return signer.Sign(url, 60*time.Second)
}

func verifyURL(r *http.Request) bool {
	url := r.URL.String()

	// The lambda function is called with a absolute URL, but the standalone
	// version is called with a relative URL. We need to make sure we have an
	// absolute URL before we can verify it.
	if !r.URL.IsAbs() {
		url = "https://" + r.Host + url
	}

	signer := surl.New([]byte(appConfig.SecretKey))
	err := signer.Verify(url)
	if err != nil {
		zerolog.Ctx(r.Context()).Error().Err(err).Msgf("Error verifying URL")
	}
	return err == nil
}
