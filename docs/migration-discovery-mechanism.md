# 数据库迁移发现机制详解

## 迁移文件发现流程

### 1. 硬编码的目录路径

迁移系统通过硬编码的路径来查找SQL文件：

```go
// 在 pkg/utils/migration.go 第40行
migrations, err := loadMigrationFiles("internal/data/migrations")
```

**配置位置：** 这个路径是硬编码在代码中的，没有配置文件。

### 2. 文件扫描过程

`loadMigrationFiles` 函数的工作流程：

```go
func loadMigrationFiles(dir string) ([]Migration, error) {
    var migrations []Migration

    // 1. 读取目录下的所有文件
    files, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    for _, file := range files {
        // 2. 只处理.sql文件
        if !strings.HasSuffix(file.Name(), ".sql") {
            continue
        }

        // 3. 跳过migrations表创建文件
        if file.Name() == "000_create_migrations_table.sql" {
            continue
        }

        // 4. 读取文件内容
        content, err := os.ReadFile(filepath.Join(dir, file.Name()))
        if err != nil {
            return nil, err
        }

        // 5. 解析文件名获取版本号和名称
        version := strings.Split(file.Name(), "_")[0]
        name := strings.TrimSuffix(strings.Join(strings.Split(file.Name(), "_")[1:], "_"), ".sql")

        // 6. 计算文件内容的MD5校验和
        hash := md5.Sum(content)
        checksum := hex.EncodeToString(hash[:])

        // 7. 创建Migration对象
        migrations = append(migrations, Migration{
            Version:  version,
            Name:     name,
            SQL:      string(content),
            Checksum: checksum,
        })
    }

    return migrations, nil
}
```

### 3. 文件命名规范

系统通过文件名解析版本号和迁移名称：

| 文件名示例 | 版本号 | 迁移名称 |
|-----------|--------|----------|
| `001_create_tasks_table.sql` | `001` | `create_tasks_table` |
| `002_create_devices_table.sql` | `002` | `create_devices_table` |
| `003_create_task_assignment_queue.sql` | `003` | `create_task_assignment_queue` |

**解析规则：**
- 版本号：文件名中第一个 `_` 之前的部分
- 迁移名称：第一个 `_` 之后到 `.sql` 之前的部分

### 4. 排序机制

发现的迁移文件会按版本号排序：

```go
// 按版本号排序
sort.Slice(migrations, func(i, j int) bool {
    return migrations[i].Version < migrations[j].Version
})
```

这确保了迁移按正确的顺序执行。

### 5. 执行条件检查

系统会检查哪些迁移需要执行：

```go
for _, migration := range migrations {
    // 检查是否已执行
    if _, exists := appliedMigrations[migration.Version]; exists {
        continue  // 跳过已执行的迁移
    }
    
    // 执行未应用的迁移
    // ...
}
```

## 配置选项

### 当前配置方式

目前迁移目录是硬编码的：
```go
migrations, err := loadMigrationFiles("internal/data/migrations")
```

### 如果要支持配置化

可以通过以下方式实现：

#### 1. 通过环境变量配置
```go
migrationsDir := os.Getenv("MIGRATIONS_DIR")
if migrationsDir == "" {
    migrationsDir = "internal/data/migrations" // 默认值
}
```

#### 2. 通过配置文件配置
```yaml
# config/config.yaml
database:
  migrations_dir: "internal/data/migrations"
```

#### 3. 通过函数参数配置
```go
func RunMigrations(db *sql.DB, migrationsDir string) error {
    // ...
}
```

## 实际的文件发现过程

让我们看看系统实际发现了哪些文件：

1. **扫描目录：** `internal/data/migrations/`
2. **找到的SQL文件：**
   - `000_create_migrations_table.sql` (跳过)
   - `001_create_tasks_table.sql` ✅
   - `002_create_devices_table.sql` ✅  
   - `003_create_task_assignment_queue.sql` ✅
   - `003_create_task_assignment_queue_enhanced.sql` ✅
   - `004_create_device_task_queue.sql` ✅

3. **解析结果：**
```
Version: 001, Name: create_tasks_table
Version: 002, Name: create_devices_table  
Version: 003, Name: create_task_assignment_queue
Version: 003, Name: create_task_assignment_queue_enhanced
Version: 004, Name: create_device_task_queue
```

4. **排序后按顺序执行**

## 总结

迁移文件的发现是通过以下机制实现的：

1. **硬编码路径：** `internal/data/migrations`
2. **文件过滤：** 只处理 `.sql` 后缀的文件
3. **命名规范：** `版本号_迁移名称.sql`
4. **自动排序：** 按版本号字符串排序
5. **去重执行：** 检查 `schema_migrations` 表避免重复执行

如果需要更改迁移文件的位置，目前需要修改代码中的硬编码路径。 