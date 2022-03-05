package verifycode

type Store interface {
	Set(key string, value string) bool
	Get(key string, clear bool) string
	Verify(key string, answer string, clear bool) bool
}
