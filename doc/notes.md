# 开发笔记

## 1. Go `range` 循环的值拷贝陷阱

`for _, entry := range slice` 中的 `entry` 是切片元素的**副本**，不是引用。

- Go 在循环开始时分配一个局部变量 `entry`，每次迭代将当前元素的值**拷贝**进去。
- `&entry` 始终是同一个地址（局部变量的地址），不是切片中原始元素的地址。
- 对 `entry` 的修改不会影响原始切片。

**错误写法**（修改的是副本，原始数据不变）：

```go
for _, entry := range workout.Entries {
    tx.QueryRow(query, ...).Scan(&entry.ID) // 写入副本，循环结束后丢失
}
```

**正确写法**（通过索引直接修改原始切片元素）：

```go
for i := range workout.Entries {
    tx.QueryRow(query, ...).Scan(&workout.Entries[i].ID) // 写入原始元素
}
```

---

## 2. 数据库迁移：`MigrateFS` vs `Migrate`

两者都调用 `goose.Up()` 执行迁移，区别在于 SQL 文件的**读取来源**。

### `MigrateFS`（生产环境 — `app.go`）

```go
store.MigrateFS(pbDb, migrations.FS, ".")
```

- 使用 `embed.FS`，SQL 文件在**编译时**被打包进二进制文件。
- 路径 `"."` 是相对于 embed 虚拟文件系统的根目录。
- 部署时只需一个二进制文件，不依赖外部 `.sql` 文件。

### `Migrate`（测试环境 — `workout_store_test.go`）

```go
Migrate(db, "../../migrations/")
```

- 直接从**磁盘文件系统**读取 SQL 文件。
- 路径是相对于测试文件所在目录（`internal/store/`）的真实磁盘路径。
- 测试始终在开发机器上运行，SQL 文件一定存在，无需嵌入。

### 对比总结

| | `MigrateFS` (app.go) | `Migrate` (test) |
|---|---|---|
| SQL 来源 | `embed.FS`（编译进二进制） | 磁盘文件系统 |
| 路径含义 | 相对于 embed 虚拟文件系统的根 | 相对于测试文件的真实磁盘路径 |
| 适用场景 | 生产 / 部署 | 本地开发 / 测试 |
