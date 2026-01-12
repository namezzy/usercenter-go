package main

import (
"context"
"fmt"
"log"
"net/http"
"time"
"user-center/internal/config"
"user-center/internal/controller"
"user-center/internal/middleware"
"user-center/internal/repository"
"user-center/internal/service"

"github.com/gin-contrib/sessions"
"github.com/gin-contrib/sessions/cookie"
"github.com/gin-gonic/gin"
"github.com/redis/go-redis/v9"
"gorm.io/driver/mysql"
"gorm.io/gorm"
)

func main() {
// 加载配置
cfg, err := config.LoadConfig("config.yaml")
if err != nil {
log.Fatalf("Failed to load config: %v", err)
}

// 连接数据库
db, err := initDB(cfg.Database)
if err != nil {
log.Fatalf("Failed to connect database: %v", err)
}

// 连接Redis
redisClient, err := initRedis(cfg.Redis)
if err != nil {
log.Fatalf("Failed to connect redis: %v", err)
}
defer redisClient.Close()

// 初始化仓库
userRepo := repository.NewUserRepository(db)

// 初始化服务
cacheService := service.NewCacheService(redisClient)
userService := service.NewUserService(userRepo, cacheService)

// 初始化控制器
userController := controller.NewUserController(userService)

// 初始化Gin
r := gin.Default()

// 配置Session
store := cookie.NewStore([]byte(cfg.Session.Secret))
store.Options(sessions.Options{
MaxAge:   cfg.Session.Timeout,
Path:     "/",
HttpOnly: true,
})
r.Use(sessions.Sessions("user-center-session", store))

// 健康检查接口
r.GET("/api/health", func(c *gin.Context) {
ctx := context.Background()

// 检查数据库
sqlDB, _ := db.DB()
dbStatus := "UP"
if err := sqlDB.Ping(); err != nil {
dbStatus = "DOWN"
}

// 检查Redis
redisStatus := "UP"
if err := cacheService.Ping(ctx); err != nil {
redisStatus = "DOWN"
}

c.JSON(http.StatusOK, gin.H{
"status":  "UP",
"message": "User Center Service is running",
"time":    time.Now().Format("2006-01-02 15:04:05"),
"components": gin.H{
"database": dbStatus,
"redis":    redisStatus,
},
})
})

// 配置路由
api := r.Group(cfg.Server.ContextPath)
{
userGroup := api.Group("/user")
{
// 公开接口
userGroup.POST("/register", userController.Register)
userGroup.POST("/login", userController.Login)

#// 需要登录

userGroup.POST("/logout", middleware.AuthMiddleware(), userController.Logout)
userGroup.GET("/current", middleware.AuthMiddleware(), userController.GetCurrentUser)

// 需要管理员权限的接口
userGroup.GET("/search", middleware.AdminAuthMiddleware(), userController.SearchUsers)
userGroup.POST("/delete", middleware.AdminAuthMiddleware(), userController.DeleteUser)
}
}

// echo
addr := fmt.Sprintf(":%d", cfg.Server.Port)
log.Printf("Server starting on %s", addr)
log.Printf("Redis connected to: %s", cfg.Redis.GetRedisAddr())
if err := r.Run(addr); err != nil {
log.Fatalf("Failed to start server: %v", err)
}
}

// initDB 初始化数据库连接
func initDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
dsn := cfg.GetDSN()

// 重试连接数据库（Docker 环境下可能需要等待 MySQL 启动）
var db *gorm.DB
var err error
maxRetries := 10

for i := 0; i < maxRetries; i++ {
db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err == nil {
break
}
log.Printf("Failed to connect database (attempt %d/%d): %v", i+1, maxRetries, err)
time.Sleep(time.Second * 3)
}

if err != nil {
return nil, fmt.Errorf("failed to connect database after %d attempts: %w", maxRetries, err)
}

sqlDB, err := db.DB()
if err != nil {
return nil, err
}

// 配置连接池
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)

log.Println("Database connected successfully")
return db, nil
}

// initRedis 初始化Redis连接
func initRedis(cfg config.RedisConfig) (*redis.Client, error) {
client := redis.NewClient(&redis.Options{
Addr:     cfg.GetRedisAddr(),
Password: cfg.Password,
DB:       cfg.DB,
})

cd /root/p_website && git commit -m "Update page title to 'Levi's Home Page'" 
ctx := context.Background()
maxRetries := 10
var err error

for i := 0; i < maxRetries; i++ {
err = client.Ping(ctx).Err()
if err == nil {
break
}
log.Printf("Failed to connect redis (attempt %d/%d): %v", i+1, maxRetries, err)
time.Sleep(time.Second * 3)
}

if err != nil {
return nil, fmt.Errorf("failed to connect redis after %d attempts: %w", maxRetries, err)
}

log.Println("Redis connected successfully")
return client, nil
}
