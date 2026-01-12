# 修复 "Unknown column 'user_account'" 错误

## 问题原因

数据库表是在修复前创建的，字段名是蛇形命名（`user_account`），但现在代码使用驼峰命名（`userAccount`）。

## 解决方案

**必须删除旧的数据卷并重新初始化数据库**

### 方法一：使用脚本（推荐）

```bash
./check_and_fix.sh
```

### 方法二：手动执行

```bash
# 1. 停止并删除所有容器和数据卷（-v 很重要！）
docker-compose down -v

# 2. 重新构建并启动
docker-compose up -d --build

# 3. 等待 30 秒让 MySQL 完全启动

# 4. 验证表结构（可选）
docker exec user-center-mysql mysql -uroot -padmin@123321 levi -e "SHOW COLUMNS FROM user;"
```

### 验证成功

执行后应该看到驼峰命名的字段：
- ✅ `userAccount` (不是 user_account)
- ✅ `avatarUrl` (不是 avatar_url)
- ✅ `userPassword` (不是 user_password)
- ✅ `planetCode` (不是 planet_code)

## 重要提示

⚠️  `-v` 参数会删除所有数据库数据！如果有重要数据，请先备份。

## 测试

重启后，访问注册接口应该可以正常工作：

```bash
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "userAccount": "testuser",
    "userPassword": "12345678",
    "checkPassword": "12345678",
    "planetCode": "12345"
  }'
```

期望返回：
```json
{
  "code": 0,
  "data": 1,
  "message": "ok"
}
```
