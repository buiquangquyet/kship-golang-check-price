package bootstrap

import (
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/infra/repo"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func BuildStorageModules() fx.Option {
	return fx.Options(
		fx.Provide(newPostgresqlDB),
		fx.Provide(newCacheRedis),
		fx.Provide(newMongoDB),
		fx.Provide(repo.NewBaseRepo),

		fx.Provide(repo.NewPriceRepo),
		fx.Provide(repo.NewShopRepo),
	)
}

func newPostgresqlDB(lc fx.Lifecycle, log *zap.SugaredLogger) *gorm.DB {
	cf := configs.Get().Postgresql
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", cf.Host,
		cf.Port, cf.User, cf.DbName, cf.SslMode, cf.Password)
	logMode := logger.Info
	if constant.IsProdEnv() {
		logMode = logger.Silent
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	//if configs.Get().Postgresql.AutoMigrate {
	//	_ = db.AutoMigrate(domain.Account{}, domain.RefreshToken{}, domain.Template{}, domain.QuotaZns{}, domain.Message{}, domain.AccountTemplateLogo{}, domain.Transaction{})
	//}
	if err != nil {
		panic(err)
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

func newMongoDB(lc fx.Lifecycle, logger *zap.SugaredLogger) *mongo.Database {
	logger.Debugf("Coming Create Storage")
	cf := configs.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.ClientOptions{}
	if configs.Get().Tracer.Enabled {
		opts.Monitor = otelmongo.NewMonitor()
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cf.Mongo.Uri), &opts)
	if err != nil {
		logger.Fatalf("connect mongo db error:[%s]", err.Error())
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		logger.Fatalf("ping mongo db error:[%s]", err.Error())
	}
	db := client.Database(cf.Mongo.DB)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Coming OnStop Storage")
			return client.Disconnect(ctx)
		},
	})
	return db
}
