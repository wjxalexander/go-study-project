package migrations

import (
	"embed"
)

// find all sql files in this migrations directory, 这通常是在讲 Go 的 embed 包，用于把 SQL 迁移文件嵌入到编译后的二进制文件中。

/**
* 将 migrations 目录下的所有 SQL 文件嵌入到编译后的二进inary file.
* 返回一个实现了 fs.FS 接口的文件系统对象。
* @return embed.FS 文件系统对象
* 这句话的意思是：“编译器，请把当前目录下所有的 .sql 文件读取出来，
* 直接以二进制数据的形式，打包进 FS 这个变量里。”
 */
//go:embed *.sql
var FS embed.FS
