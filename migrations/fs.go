package migrations

import (
	"embed"
)

// find all sql files in this migrations directory, 这通常是在讲 Go 的 embed 包，用于把 SQL 迁移文件嵌入到编译后的二进制文件中。

//go:embed *.sql
var FS embed.FS
