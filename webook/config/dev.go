//go:build !k8s

package config

var Config = config{
	DB:    DBConfig{DNS: "root:root@tcp(localhost:13316)/webook"},
	Redis: RedisConfig{Addr: "localhost:6379"},
}
