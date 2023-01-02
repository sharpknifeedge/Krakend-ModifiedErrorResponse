package consts

const (
	HTTPUserID = "X-User-ID"

	//Deprecated:HTTPROleHeader use `HTTPRoleHeader` instead
	HTTPROleHeader = "X-User-Role"
	HTTPRoleHeader = "X-User-Role"

	HTTPTokenRole     = "X-Token-Role"
	HTTPTokenUserID   = "X-Token-User-ID"
	HTTPRequestID     = "X-Request-Id"
	HTTPXForwardedFor = "X-Forwarded-For"
)
