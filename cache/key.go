package cache

func Get(key string) (string, error) {
	return rc.Get(key).Result()
}

func Set(key, value string) (string, error) {
	return rc.Set(key, value, 0).Result()
}
