package spa

import (
	"fmt"
	"net/http"
	"strings"
)

type HSTSMiddleware struct {
	h                 http.Handler
	seconds           uint64
	includeSubdomains bool
	preload           bool
}

func NewHSTSMiddleware(next http.Handler, hstsSeconds uint64, includeSubdomains bool, preload bool) (*HSTSMiddleware, error) {
	return &HSTSMiddleware{next, hstsSeconds, includeSubdomains, preload}, nil
}

func (o *HSTSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if GetRequestScheme(r) == "https" {
		stsParameters := []string{fmt.Sprintf("max-age=%d", o.seconds)}
		if o.includeSubdomains {
			stsParameters = append(stsParameters, "includeSubDomains")
		}
		if o.preload {
			stsParameters = append(stsParameters, "preload")
		}
		stsString := strings.Join(stsParameters[:], "; ")
		w.Header().Set("Strict-Transport-Security", stsString)
	}
	o.h.ServeHTTP(w, r)
}
