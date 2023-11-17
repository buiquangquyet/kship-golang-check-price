package bootstrap

import (
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/infra/decorators"
	"check-price/src/infra/repo"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func BuildStorageModules() fx.Option {
	return fx.Options(
		fx.Provide(newMysqlDB),
		fx.Provide(newCacheRedis),

		fx.Provide(repo.NewBaseRepo),
		fx.Provide(repo.NewCityRepo),
		fx.Provide(repo.NewClientRepo),
		fx.Provide(repo.NewDistrictRepo),
		fx.Provide(repo.NewServiceRepo),
		fx.Provide(repo.NewSettingShopRepo),
		fx.Provide(repo.NewShopRepo),
		fx.Provide(repo.NewWardRepo),
		fx.Provide(repo.NewConfigCodT0Repo),
		fx.Provide(repo.NewSettingRepo),
		fx.Provide(repo.NewShopCodT0Repo),
		fx.Provide(repo.NewShopLevelRepo),

		fx.Provide(decorators.NewBaseDecorator),
		fx.Provide(decorators.NewCityRepoDecorator),
		fx.Provide(decorators.NewClientRepoDecorator),
		fx.Provide(decorators.NewDistrictRepoDecorator),
		fx.Provide(decorators.NewServiceRepoDecorator),
		fx.Provide(decorators.NewSettingShopRepoDecorator),
		fx.Provide(decorators.NewShopRepoDecorator),
		fx.Provide(decorators.NewWardRepoDecorator),
		fx.Provide(decorators.NewConfigCodT0RepoDecorator),
		fx.Provide(decorators.NewSettingRepoDecorator),
		fx.Provide(decorators.NewShopCodT0RepoDecorator),
		fx.Provide(decorators.NewShopLevelRepoDecorator),
	)
}

func newMysqlDB(lc fx.Lifecycle, log *zap.SugaredLogger) *gorm.DB {
	cf := configs.Get().Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cf.User, cf.Password,
		cf.Host, cf.Port, cf.DbName)
	logMode := logger.Info
	if constant.IsProdEnv() {
		logMode = logger.Silent
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		log.Fatal(err)
	}
	if configs.Get().Tracer.Enabled {
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			panic(err)
		}
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Debug("Coming OnStop Storage")
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})
	return db
}

func newCacheRedis() redis.UniversalClient {
	cf := configs.Get().Redis
	hosts := cf.Hosts
	var client redis.UniversalClient
	isClusterMode := len(hosts) > 1
	if isClusterMode {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    hosts,
			Username: cf.Username,
			Password: cf.Password,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     hosts[0],
			Username: cf.Username,
			Password: cf.Password,
		})
	}

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.GetLogger().GetZap().Fatalf("ping redis error, err:[%s]", err.Error())
	}
	return client
}
