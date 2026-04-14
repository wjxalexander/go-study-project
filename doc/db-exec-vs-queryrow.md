# `database/sql`：`Exec` 与 `QueryRow` 的区别

在 Go 标准库 `database/sql` 中，`Exec` 和 `QueryRow` 都通过 `*sql.DB`（或 `*sql.Tx`）执行 SQL，但**适用场景和返回值完全不同**。

## 一句话区分

| 方法 | 典型用途 | 返回什么 |
|------|----------|----------|
| **`Exec`** | 不返回结果集的语句（写操作、DDL） | `sql.Result`（受影响行数、`LastInsertId` 等） |
| **`QueryRow`** | 预期**恰好一行**的查询 | `*sql.Row`，用 `Scan` 取列 |

---

## `Exec`

用于执行 **INSERT、UPDATE、DELETE**以及 **CREATE TABLE** 等**不返回行集**的语句。

- 成功时得到 `sql.Result`，常用：
  - `RowsAffected()`：受影响行数
  - `LastInsertId()`：自增主键（依赖驱动与表定义，PostgreSQL 等场景常不可靠，习惯用 `RETURNING` + `QueryRow`）
- **不要**用 `Exec` 去读 `SELECT` 的结果；若误用，驱动可能忽略结果集或表现不符合预期。

```go
res, err := db.Exec(`UPDATE users SET name = $1 WHERE id = $2`, name, id)
if err != nil {
    return err
}
n, err := res.RowsAffected()
// ...
```

---

## `QueryRow`

用于 **SELECT**（或带 `RETURNING` 的 INSERT/UPDATE）且你**只关心一行**的情况。

- 内部仍会向数据库发查询；返回的 `*sql.Row` 在调用 `Scan` 时才真正取数。
- 若 **没有行**：`Scan` 返回 **`sql.ErrNoRows`**（这是 `QueryRow` 的惯用错误处理方式）。
- 若 **多行**：`QueryRow` **只取第一行**，其余行会被丢弃（容易埋 bug），多行应使用 **`Query`** +循环 `Rows.Next()`。

```go
var name string
err := db.QueryRow(`SELECT name FROM users WHERE id = $1`, id).Scan(&name)
if err == sql.ErrNoRows {
    // 未找到
}
if err != nil {
    return err
}
```

---

## 和 `Query` 的关系（顺带）

- **`Query`**：返回 `*sql.Rows`，适合 **零行、一行或多行**；要手动 `defer rows.Close()`，并用 `Next()`遍历。
- **`QueryRow`**：等价于「只读第一行的便捷封装」，语义是「我期望最多一行」。

---

## 实务建议

1. **写库 / DDL → `Exec`**；**读单行 → `QueryRow` + `Scan`**；**读多行 → `Query`**。
2. PostgreSQL 拿插入后的主键或字段，常用 **`INSERT ... RETURNING id` + `QueryRow`**，而不是依赖 `LastInsertId()`。
3. 需要事务时，在 `*sql.Tx` 上调用同名方法，语义与 `*sql.DB` 一致。
