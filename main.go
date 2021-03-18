/**
 * @Author: Bugzheng
 * @Description:
 * @File:  main.go
 * @Version: 1.0.0
 * @Date: 2021/02/01 3:12 下午
 */
package main

import (
	"JoGo/app/api"
	"JoGo/app/model"
	"JoGo/boot"
	"JoGo/pkg/cmd"
	"JoGo/pkg/middlewares"
	"JoGo/route"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "JoGo"
	app.Usage = " ////////////////////////////////////////////////////////////////////\n//                          _ooOoo_                               //\n//                         o8888888o                              //\n//                         88\" . \"88                              //\n//                         (| ^_^ |)                              //\n//                         O\\  =  /O                              //\n//                      ____/`---'\\____                           //\n//                    .'  \\\\|     |//  `.                         //\n//                   /  \\\\|||  :  |||//  \\                        //\n//                  /  _||||| -:- |||||-  \\                       //\n//                  |   | \\\\\\  -  /// |   |                       //\n//                  | \\_|  ''\\---/''  |   |                       //\n//                  \\  .-\\__  `-`  ___/-. /                       //\n//                ___`. .'  /--.--\\  `. . ___                     //\n//              .\"\" '<  `.___\\_<|>_/___.'  >'\"\".                  //\n//            | | :  `- \\`.;`\\ _ /`;.`/ - ` : | |                 //\n//            \\  \\ `-.   \\_ __\\ /__ _/   .-` /  /                 //\n//      ========`-.____`-.___\\_____/___.-`____.-'========         //\n//                           `=---='                              //\n//      ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^        //\n//                  佛祖保佑       永不宕机     永无BUG              //\n////////////////////////////////////////////////////////////////////\n    "
	app.Commands = cmd.Commands
	app.Action = func(c *cli.Context) error {
		boot.ZapLogger.Info("项目版本V1.0启动成功")
		appInit()
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("app running error", err)
		boot.ZapLogger.Fatal("app running error", zap.Error(err))
	}
}

func appInit() {
	if env := os.Getenv("DEPLOY_ENV"); env == "dev" {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	//初始化路由
	route.RegisterApp(r)
	//数据库表的自动迁移
	r.Use(middlewares.Recovery())
	model.Migration()
	if err := api.TransInit("zh"); err != nil {
		boot.ZapLogger.Error("翻译中文初始化失败init trans failed, err:%v\n", zap.Error(err))
		return
	}
	// go tool pprof （go 自带的性能分析工具）
	pprof.Register(r)
	// Graceful restart & zero downtime deploy for Go servers.
	// Use `kill -USR2 pid` to restart(使用 `kill -USR2 进程号热重启新的进程dddd，优雅退出旧的服务`).
	if err := gracehttp.Serve(
		&http.Server{
			Addr:         fmt.Sprintf(":%d", 8998),
			Handler:      r,
			IdleTimeout:  10 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}); err != nil {
		boot.ZapLogger.Fatal("serving error", zap.Error(err))
	}
}
