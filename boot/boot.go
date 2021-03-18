//boot 项目启动文件
/**
 * @Author: BugZheng
 * @Description:初始化的
 * @File:  boot
 * @Version: 1.0.0
 * @Date: 2021/02/22 5:54 下午
 */
package boot

import (
	"JoGo/pkg/conf"
	"JoGo/pkg/logger"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

//Config 配置文件的结构体(人工智能：将yaml 转成json  再将json 转成 结构体)
//小工具1：https://www.bejson.com/json/json2yaml/
//小工具2：https://mholt.github.io/json-to-go/
type Config struct {
	LogsPath       string `json:"LogsPath"`
	DataBaseConfig struct {
		DBHOST     string `json:"DB_HOST"`
		DBPORT     int    `json:"DB_PORT"`
		DBUSERNAME string `json:"DB_USERNAME"`
		DBPASSWORD int    `json:"DB_PASSWORD"`
		DBNAME     string `json:"DB_NAME"`
	} `json:"DataBaseConfig"`
}

// DB 数据库链接单例
var DB *gorm.DB

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client

var ZapLogger *zap.Logger

// Init 初始化配置项
func init() {
	conf.ConfigInit()
	InitDatabase(os.Getenv("MYSQL_DSN"))
	InitRedis()
	InitLogger()
}

//InitDatabase 在中间件中初始化mysql链接
func InitDatabase(connString string) {
	//connString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
	//	GetConfig().DataBaseConfig.DBHOST,
	//	GetConfig().DataBaseConfig.DBPASSWORD,
	//	GetConfig().DataBaseConfig.DBHOST,
	//	GetConfig().DataBaseConfig.DBNAME)
	connString = "root:123456@tcp(120.78.152.4:3306)/ego_db?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", connString)
	db.LogMode(true)
	// Error
	if err != nil {
		ZapLogger.Error("连接数据库不成功", zap.Error(err))
		os.Exit(0)
	}
	//设置连接池
	//空闲
	db.DB().SetMaxIdleConns(50)
	//打开
	db.DB().SetMaxOpenConns(100)
	//超时
	db.DB().SetConnMaxLifetime(time.Second * 30)

	DB = db
}

//InitRedis 在中间件中初始化redis链接
func InitRedis() {
	db, _ := strconv.ParseUint(os.Getenv("REDIS_DB"), 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:       os.Getenv("REDIS_ADDR"),
		Password:   os.Getenv("REDIS_PW"),
		DB:         int(db),
		MaxRetries: 1,
	})

	_, err := client.Ping().Result()

	if err != nil {
		ZapLogger.Error("连接Redis不成功", zap.Error(err))
		os.Exit(0)
	}

	RedisClient = client
}

/**
 * 获取日志
 * filePath 日志文件路径
 * level 日志级别
 * maxSize 每个日志文件保存的最大尺寸 单位：M
 * maxBackups 日志文件最多保存多少个备份
 * maxAge 文件最多保存多少天
 * compress 是否压缩
 * serviceName 服务名
 */
//InitLogger 初始化
func InitLogger() *zap.Logger {
	//Debug 的时候日志在控制台输出
	ZapLogger = logger.NewLogger(&logger.LogConfig{
		Path:       "logs/app.log",
		MaxSize:    10,
		MaxBackups: 1,
		MaxAge:     10,
		Compress:   true,
	}, conf.GetBool("DEBUG"))
	ZapLogger.Info("log 初始化成功")
	return ZapLogger
}
