# JWT Authentication Middleware

## 概述

JWT middleware 提供員工用戶的 JWT 令牌驗證功能，確保 API 安全訪問。

## 使用方式

### 基本使用

```go
package main

import (
  "github.com/gin-gonic/gin"
  "github.com/tkoleo84119/nail-salon-backend/internal/config"
  "github.com/tkoleo84119/nail-salon-backend/internal/middleware"
  "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
)

func main() {
  // 初始化配置和資料庫
  cfg := config.Load()
  db := dbgen.New(dbConnection)

  router := gin.Default()

  // 使用JWT認證中介層
  protected := router.Group("/api/v1")
  protected.Use(middleware.JWTAuth(cfg, db))

  // 保護的路由
  protected.GET("/profile", getProfile)
  protected.POST("/logout", logout)
}
```

### 獲取用戶上下文

在受保護的路由處理器中獲取當前用戶信息：

```go
func getProfile(c *gin.Context) {
  // 從context中獲取員工信息
  staff, exists := middleware.GetStaffFromContext(c)
  if !exists {
    c.JSON(http.StatusUnauthorized, common.ErrorResponse("獲取員工資訊失敗", nil))
    return
  }

  // 使用員工信息
  response := ProfileResponse{
    UserID:    staff.UserID,
    Username:  staff.Username,
    Role:      staff.Role,
    StoreList: staff.StoreList,
  }

  c.JSON(http.StatusOK, common.SuccessResponse(response))
}
```

## 錯誤回應

middleware 會根據不同的認證失敗情況返回相應的錯誤訊息：

### 缺少認證令牌
```json
{
  "message": "認證失敗",
  "errors": {
    "token": "access_token 缺失"
  }
}
```

### 令牌格式錯誤
```json
{
  "message": "認證失敗",
  "errors": {
    "token": "access_token 格式錯誤"
  }
}
```

### 令牌無效或過期
```json
{
  "message": "認證失敗",
  "errors": {
    "token": "access_token 無效或已過期"
  }
}
```

### 員工認證失敗
```json
{
  "message": "認證失敗",
  "errors": {
    "token": "員工認證失敗"
  }
}
```

## 技術細節

### 令牌驗證流程

1. 檢查Authorization標頭是否存在
2. 驗證令牌格式（Bearer token）
3. 使用JWT工具驗證令牌簽名和有效期
4. 從資料庫驗證員工用戶是否存在且為啟用中（`is_active = true`）
5. 將員工上下文存儲到 `Gin context` 中

### 擴展性考慮

當前實作僅支援員工用戶認證。未來可擴展支援客戶認證。