.PHONY: mock
mock:
	@mockgen -source=./webook/internal/service/user.go -package=svcmocks -destination=./webook/internal/service/mocks/user.mock.go
	@mockgen -source=./webook/internal/service/code.go -package=svcmocks -destination=./webook/internal/service/mocks/code.mock.go
	@mockgen -source=./webook/internal/service/sms/types.go -package=smsmocks -destination=./webook/internal/service/sms/mocks/sms.mock.go
	@mockgen -source=./webook/internal/repository/code.go -package=repmocks -destination=./webook/internal/repository/mocks/code.mock.go
	@mockgen -source=./webook/internal/repository/user.go -package=repmocks -destination=./webook/internal/repository/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/user.go -package=repmocks -destination=./webook/internal/repository/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/cache/user.go -package=cachemocks -destination=./webook/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/cache/code.go -package=cachemocks -destination=./webook/internal/repository/cache/mocks/code.mock.go
	@mockgen -source=./webook/internal/repository/dao/user.go -package=daomocks -destination=./webook/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./webook/pkg/limiter/types.go -package=limitmocks -destination=./webook//pkg/limiter/mocks/limiter.mock.go
	@go mod tidy