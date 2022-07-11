package rester

import "net/http"

/* This function will check whether we need to try calling again to the hitting server.
It calls back in case of
- 	HTTP 102 (Status Processing)
-	HTTP ERR 408 (Request Timeout)
-	HTTP ERR 500 (Internal Server Error)
-	HTTP ERR 502 (Bad Gateway)
-	HTTP ERR 503 (Service Unavailable)
-	HTTP ERR 504 (Gateway Timeout)
*/
func HttpCodeForRetry(httpStatusCode int) bool {
	switch httpStatusCode {
	case
		http.StatusRequestTimeout,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusProcessing:
		return true
	}
	return false
}
