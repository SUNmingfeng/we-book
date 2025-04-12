//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DNS: "root:root@tcp(webook-mysql-service:3308)/mysql",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:6379",
	},
}
