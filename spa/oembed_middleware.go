package spa

import (
	"fmt"
	"net/http"
	"net/url"
)

type OembedMiddleware struct {
	h   http.Handler
	URL *url.URL
}

func NewOembedMiddleware(next http.Handler, urlStr string) (*OembedMiddleware, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return &OembedMiddleware{next, u}, nil
}

func (o *OembedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Prepare the target url
	targetUrl := cloneUrl(r.URL)
	targetUrl.Host = r.Host

	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if r.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}
	targetUrl.Scheme = scheme

	// Build the oEmbed URL
	u := cloneUrl(o.URL)
	q := u.Query()
	q.Set("url", targetUrl.String())
	u.RawQuery = q.Encode()

	// Set the Link HTTP header
	v := fmt.Sprintf(`<%s>; rel="alternate"; type="application/json+oembed"`, u.String())
	w.Header().Set("Link", v)
	o.h.ServeHTTP(w, r)
}

func cloneUrl(u *url.URL) *url.URL {
	u2 := *u
	return &u2
}
