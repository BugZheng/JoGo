/**
 * @Author: zhengweizhao
 * @Description:
 * @File:  codeGenerator
 * @Version: 1.0.0
 * @Date: 2021/03/18 5:31 下午
 */
package cmd

import (
	"JoGo/boot"
	"JoGo/pkg/conf"
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
	"strings"
	"text/template"
)

const (
	//dbName             =  conf.Get("DB_NAME")
	packageName = "model"
)

//生成代码
func Code(tableName string) {
	ctx := context.Background()

	info := selectTableColumn(ctx, tableName)
	//创建目录
	os.MkdirAll("app/model", os.ModePerm)
	structFileName := "app/model/" + info["structName"].(string) + ".go"
	structFile, _ := os.Create(structFileName)
	defer func() {
		structFile.Close()
	}()

	structTemplate, err1 := template.ParseFiles("./pkg/cmd/templates/struct.txt")
	if err1 != nil {
		fmt.Println(err1)
	}
	structTemplate.Execute(structFile, info)

	//serviceTemplate, err2 := template.ParseFiles("./templates/service.txt")
	//if err2 != nil {
	//	fmt.Println(err2)
	//}
	//serviceTemplate.Execute(serviceFile, info)

}

//获取所有的表名
func selectAllTable() []string {
	tableNames := []string{}
	boot.DB.Raw("select table_name from information_schema.TABLES where  TABLE_SCHEMA =?", conf.Get("DB_NAME")).Scan(&tableNames)
	return tableNames
}

//根据表名查询字段信息和主键名称
func selectTableColumn(ctx context.Context, tableName string) map[string]interface{} {

	info := make(map[string]interface{})

	tableComment := ""
	rows, err := boot.DB.Raw("select table_comment from information_schema.TABLES where  TABLE_SCHEMA =? and TABLE_Name=? ", conf.Get("DB_NAME"), tableName).Rows()
	if err != nil {
		boot.ZapLogger.Error("查询数据库失败", zap.Error(err))
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&tableComment)
	}

	//查找主键
	pkName := ""
	rows2, err := boot.DB.Raw("SELECT column_name FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA=? and  table_name=? AND constraint_name=?", conf.Get("DB_NAME"), tableName, "PRIMARY").Rows()
	if err != nil {
		boot.ZapLogger.Error("查询数据库失败", zap.Error(err))
		return nil
	}
	defer rows2.Close()
	for rows2.Next() {
		rows2.Scan(&pkName)
	}
	maps := []map[string]interface{}{}
	if len(maps) == 0 {
		boot.ZapLogger.Info("表不存在")
		fmt.Println("表不存在")
		return nil
	}
	// select * from information_schema.COLUMNS where table_schema ='readygo' and table_name='t_user';
	rows3, err := boot.DB.Raw("select COLUMN_NAME,DATA_TYPE,IS_NULLABLE,COLUMN_COMMENT from information_schema.COLUMNS where  TABLE_SCHEMA =? and TABLE_NAME=? and COLUMN_NAME not like ?  order by ORDINAL_POSITION asc", conf.Get("DB_NAME"), tableName, "bak%").Rows()
	defer rows3.Close()
	columnName := ""
	dataType := ""
	isNullAble := ""
	columnComment := ""
	for rows3.Next() {
		rows3.Scan(&columnName, &dataType, &isNullAble, &columnComment)
		maps = append(maps, map[string]interface{}{
			"DATA_TYPE":      dataType,
			"COLUMN_NAME":    columnName,
			"IS_NULLABLE":    isNullAble,
			"COLUMN_COMMENT": columnComment,
		})
	}
	for _, m := range maps {
		dataType := m["DATA_TYPE"].(string)
		dataType = strings.ToUpper(dataType)

		nullable := m["IS_NULLABLE"].(string)
		nullable = strings.ToUpper(nullable)

		if dataType == "VARCHAR" || dataType == "NVARCHAR" || dataType == "TEXT" || dataType == "LONGTEXT" {
			//if nullable == "YES" {
			//	dataType = "sql.NullString"
			//} else {
			dataType = "string"
			//}

		} else if dataType == "DATETIME" || dataType == "TIMESTAMP" {
			//if nullable == "YES" {
			//	dataType = "sql.NullTime"
			//} else {
			dataType = "time.Time"
			//}

		} else if dataType == "INT" {
			//if nullable == "YES" {
			//	dataType = "sql.NullInt32"
			//} else {
			dataType = "int"
			//}

		} else if dataType == "BIGINT" {
			//if nullable == "YES" {
			//	dataType = "sql.NullInt64"
			//} else {
			dataType = "int64"
			//}
		} else if dataType == "SMALLINT" {
			dataType = "int32"
		} else if dataType == "FLOAT" {
			//if nullable == "YES" {
			//	dataType = "sql.NullFloat64"
			//} else {
			dataType = "float32"
			//}

		} else if dataType == "DOUBLE" {
			//if nullable == "YES" {
			//	dataType = "sql.NullFloat64"
			//} else {
			dataType = "float64"
			//}

		} else if dataType == "DECIMAL" {
			dataType = "decimal.Decimal"
		} else if dataType == "TINYINT" {
			dataType = "int8"
		}
		m["DATA_TYPE"] = dataType
		fieldName := camelCaseName(m["COLUMN_NAME"].(string))
		m["field"] = fieldName

		//设置主键的struct属性名称
		if m["COLUMN_NAME"].(string) == pkName {
			info["pkField"] = fieldName
		}

	}

	info["columns"] = maps
	info["pkName"] = pkName
	info["tableName"] = tableName
	structName := tableName
	if strings.HasPrefix(structName, "t_") {
		structName = structName[2:]
	}
	structName = camelCaseName(structName) + "Struct"
	info["structName"] = structName
	info["pname"] = firstToLower(structName)
	info["packageName"] = packageName
	info["tableComment"] = tableComment
	return info
}

//首字母大写
func firstToUpper(str string) string {
	str = strings.ToUpper(string(str[0:1])) + string(str[1:])
	return str
}

//首字母小写
func firstToLower(str string) string {
	str = strings.ToLower(string(str[0:1])) + string(str[1:])
	return str
}

//驼峰
func camelCaseName(name string) string {
	names := strings.Split(name, "_")
	structName := ""
	for _, name := range names {
		structName = structName + firstToUpper(name)
	}

	return structName

}
