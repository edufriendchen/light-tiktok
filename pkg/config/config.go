package config

type Config interface {
	GetInt(key string) int64
	GetString(key string) string
	GetBool(key string) bool
	GetFloat(key string) float64
	GetBinary(key string) []byte
	GetArray(key string) []string
	GetMap(key string) map[string]string
}
