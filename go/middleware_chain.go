package mode

import "net/http"

func Chain(
	handler http.HandlerFunc,
	name string,
	auth func(http.HandlerFunc) http.HandlerFunc,
) http.HandlerFunc {

	return MetricsMiddleware(auth(handler), name)
}
