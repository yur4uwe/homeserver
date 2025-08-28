package constants

var Rootd = ""

type CtxKey string

const (
	RequesterKey CtxKey = "requester"
	EntryKey     CtxKey = "entry"
)

const (
	AuthCookieName = "authTokenHS"
	SIDQueryKey    = "sid"
	ExpQueryKey    = "exp"
)
