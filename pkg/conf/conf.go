/**
 * @Author: BugZheng
 * @Description:读取配置文件的值
 * @File:  conf
 * @Version: 1.0.0
 * @Date: 2021/02/24 11:00 上午
 */
// Package conf 提供最基础的配置加载功能
package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var path string
var files map[string]*Conf

var (
	// Hostname 主机名
	Hostname = "localhost"
	// AppID 获取 APP_ID
	AppID = "localapp"
	// IsDevEnv 开发环境标志
	IsDevEnv = false
	// IsTestEnv 测试环境标志
	IsTestEnv = false
	// IsProdEnv 生产环境标志
	IsProdEnv = false
	// Env 运行环境
	Env = "dev"
	// Zone 服务区域
	Zone = "sh001"
)

func ConfigInit() {
	Hostname, _ = os.Hostname()
	if appID := os.Getenv("APP_ID"); appID != "" {
		AppID = appID
	}
	if env := os.Getenv("DEPLOY_ENV"); env != "" {
		Env = env
	}
	if zone := os.Getenv("ZONE"); zone != "" {
		Zone = zone
	}
	switch Env {
	case "prod", "pre":
		IsProdEnv = true
	case "test":
		IsTestEnv = true
	default:
		IsDevEnv = true
	}

	path = os.Getenv("CONF_PATH")

	if path == "" {
		//var err error
		//if path, err = os.Getwd(); err != nil {
		//	panic(err)
		//}
		path = "config"
	}

	fs, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	files = make(map[string]*Conf, len(fs))

	for _, f := range fs {
		fmt.Println(f.Name())
		if !strings.HasSuffix(f.Name(), ".toml") {
			continue
		}
		v := viper.New()
		v.SetConfigFile(path + "/" + f.Name())
		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}
		v.AutomaticEnv()
		name := strings.TrimSuffix(f.Name(), ".toml")
		files[name] = &Conf{v}
	}
}

type Conf struct {
	viper *viper.Viper
}

// GetFloat64 获取浮点数配置
func GetFloat64(key string) float64 { return File().GetFloat64(key) }
func (c *Conf) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

// Get 获取字符串配置
func Get(key string) string { return File().Get(key) }
func (c *Conf) Get(key string) string {
	return c.viper.GetString(key)
}

// GetStrings 获取字符串列表
func GetStrings(key string) (s []string) { return File().GetStrings(key) }
func (c *Conf) GetStrings(key string) (s []string) {
	value := Get(key)
	if value == "" {
		return
	}

	for _, v := range strings.Split(value, ",") {
		s = append(s, v)
	}
	return
}

// GetInt32s 获取数字列表
// 1,2,3 => []int32{1,2,3}
func GetInt32s(key string) (s []int32, err error) { return File().GetInt32s(key) }
func (c *Conf) GetInt32s(key string) (s []int32, err error) {
	s64, err := GetInt64s(key)
	for _, v := range s64 {
		s = append(s, int32(v))
	}
	return
}

// GetInt64s 获取数字列表
func GetInt64s(key string) (s []int64, err error) { return File().GetInt64s(key) }
func (c *Conf) GetInt64s(key string) (s []int64, err error) {
	value := Get(key)
	if value == "" {
		return
	}

	var i int64
	for _, v := range strings.Split(value, ",") {
		i, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		s = append(s, i)
	}
	return
}

// GetInt 获取整数配置
func GetInt(key string) int { return File().GetInt(key) }
func (c *Conf) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetInt32 获取 int32 配置
func GetInt32(key string) int32 { return File().GetInt32(key) }
func (c *Conf) GetInt32(key string) int32 {
	return c.viper.GetInt32(key)
}

// GetInt64 获取 int64 配置
func GetInt64(key string) int64 { return File().GetInt64(key) }
func (c *Conf) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

// GetDuration 获取时间配置
func GetDuration(key string) time.Duration { return File().GetDuration(key) }
func (c *Conf) GetDuration(key string) time.Duration {
	return c.viper.GetDuration(key)
}

// GetTime 查询时间配置
// 默认时间格式为 "2006-01-02 15:04:05"，conf.GetTime("FOO_BEGIN")
// 如果需要指定时间格式，则可以多传一个参数，conf.GetString("FOO_BEGIN", "2006")
//
// 配置不存在或时间格式错误返回**空时间对象**
// 使用本地时区
func GetTime(key string, args ...string) time.Time { return File().GetTime(key, args...) }
func (c *Conf) GetTime(key string, args ...string) time.Time {
	fmt := "2006-01-02 15:04:05"
	if len(args) == 1 {
		fmt = args[0]
	}

	t, _ := time.ParseInLocation(fmt, c.viper.GetString(key), time.Local)
	return t
}

// GetBool 获取配置布尔配置
func GetBool(key string) bool { return File().GetBool(key) }
func (c *Conf) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// Set 设置配置，仅用于测试
func Set(key string, value string) { File().Set(key, value) }
func (c *Conf) Set(key string, value string) {
	c.viper.Set(key, value)
}

// File 根据文件名获取对应配置对象
// 目前仅支持 toml 文件，不用传扩展名
// 如果要读取 foo.toml 配置，可以 File("foo").Get("bar")
func File() *Conf {
	//先去系统环境系统变量的配置文件找
	//找不到再去环境变量(app.conf)尝试找
	//根据环境文件的配置（默认dev）
	if IsProdEnv == true {
		return files["prod"]
	}
	if IsTestEnv == true {
		return files["test"]
	}
	if IsDevEnv == true {
		appConfRunMode := files["app"].viper.GetString("RUNMODE")
		if appConfRunMode != "" {
			return files[appConfRunMode]
		}

	}
	return files["dev"]
}
