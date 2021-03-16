/**
 * @Author: BugZheng
 * @Description:数据表的迁移
 * @File:  autoMigrate
 * @Version: 1.0.0
 * @Date: 2021/02/22 7:38 下午
 */
package model

import "JoGo/boot"

//执行数据迁移
func Migration() {
	// 自动迁移模式
	boot.DB.AutoMigrate(&User{})
}
