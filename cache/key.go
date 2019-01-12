package cache

// Get returns a value for a key.
func Get(key string) (string, error) {
	return rc.Get(key).Result()
}

// Set sets a value for a key.
func Set(key, value string) (string, error) {
	return rc.Set(key, value, 0).Result()
}
