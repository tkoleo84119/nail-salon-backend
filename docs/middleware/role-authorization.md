# Role Authorization Middleware

## 概述

Role Authorization middleware 提供員工用戶的角色權限驗證功能，確保只有具備適當權限的用戶才能訪問特定的 API 端點。此 middleware 需要與 JWT Authentication middleware 搭配使用。

## 使用方式

### 基本使用

```go
package main

import (
  "github.com/gin-gonic/gin"
  "github.com/tkoleo84119/nail-salon-backend/internal/config"
  "github.com/tkoleo84119/nail-salon-backend/internal/middleware"
  "github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
  "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
)

func main() {
  // 初始化配置和資料庫
  cfg := config.Load()
  db := dbgen.New(dbConnection)

  router := gin.Default()

  // 使用 JWT 認證和角色授權中介層
  adminRoutes := router.Group("/api/v1/admin")
  adminRoutes.Use(middleware.JWTAuth(cfg, db))
  adminRoutes.Use(middleware.RequireAdminRoles())

  // 只有 SUPER_ADMIN 和 ADMIN 可以訪問
  adminRoutes.GET("/settings", getAdminSettings)
  adminRoutes.POST("/users", createUser)

  // 使用自定義角色組合
  managerRoutes := router.Group("/api/v1/manager")
  managerRoutes.Use(middleware.JWTAuth(cfg, db))
  managerRoutes.Use(middleware.RequireRoles(staff.RoleSuperAdmin, staff.RoleAdmin, staff.RoleManager))

  managerRoutes.GET("/reports", getReports)
}
```

### 角色權限函數

#### RequireRoles(allowedRoles ...string)
核心授權函數，可傳入多個允許的角色：

```go
// 只允許 SUPER_ADMIN 和 ADMIN
router.GET("/sensitive-data",
    middleware.JWTAuth(cfg, db),
    middleware.RequireRoles(staff.RoleSuperAdmin, staff.RoleAdmin),
    handler.GetSensitiveData,
)

// 允許 MANAGER 和 STYLIST
router.POST("/booking/update",
    middleware.JWTAuth(cfg, db),
    middleware.RequireRoles(staff.RoleManager, staff.RoleStylist),
    handler.UpdateBooking,
)
```

#### 便利函數

**RequireSuperAdmin()** - 僅允許超級管理員
```go
router.DELETE("/system/reset",
    middleware.JWTAuth(cfg, db),
    middleware.RequireSuperAdmin(),
    handler.SystemReset,
)
```

**RequireAdminRoles()** - 允許管理員級別（SUPER_ADMIN, ADMIN）
```go
router.GET("/admin/reports",
    middleware.JWTAuth(cfg, db),
    middleware.RequireAdminRoles(),
    handler.GetAdminReports,
)
```

**RequireManagerOrAbove()** - 允許經理級別以上（SUPER_ADMIN, ADMIN, MANAGER）
```go
router.POST("/staff/schedule",
    middleware.JWTAuth(cfg, db),
    middleware.RequireManagerOrAbove(),
    handler.CreateStaffSchedule,
)
```

**RequireAnyStaffRole()** - 允許任何員工角色
```go
router.GET("/profile",
    middleware.JWTAuth(cfg, db),
    middleware.RequireAnyStaffRole(),
    handler.GetProfile,
)
```

## 角色常數

使用 `internal/model/staff/role.go` 中定義的角色常數：

```go
staff.RoleSuperAdmin  // "SUPER_ADMIN" - 超級管理員
staff.RoleAdmin       // "ADMIN"       - 管理員
staff.RoleManager     // "MANAGER"     - 經理
staff.RoleStylist     // "STYLIST"     - 美甲師
```

## 中介層順序

**重要：** 角色授權 middleware 必須在 JWT 認證 middleware 之後使用：

```go
// ✅ 正確順序
router.GET("/protected",
    middleware.JWTAuth(cfg, db),        // 第一步：身份認證
    middleware.RequireRoles(...),       // 第二步：權限授權
    handler.ProtectedEndpoint,
)

// ❌ 錯誤順序 - RequireRoles 需要已認證的用戶上下文
router.GET("/protected",
    middleware.RequireRoles(...),       // 這會失敗
    middleware.JWTAuth(cfg, db),
    handler.ProtectedEndpoint,
)
```

## 錯誤回應

middleware 會根據不同的授權失敗情況返回相應的錯誤訊息：

### 未找到使用者認證資訊（401）
當用戶上下文不存在時（通常是 JWT middleware 未執行或失敗）：
```json
{
  "message": "未找到使用者認證資訊"
}
```

### 權限不足（403）
當用戶角色不在允許的角色列表中時：
```json
{
  "message": "權限不足，無法執行此操作"
}
```

## 實際應用範例

### 多層級權限控制
```go
// 系統管理相關 - 僅超級管理員
systemRoutes := router.Group("/api/v1/system")
systemRoutes.Use(middleware.JWTAuth(cfg, db))
systemRoutes.Use(middleware.RequireSuperAdmin())

// 店面管理相關 - 管理員級別以上
storeRoutes := router.Group("/api/v1/store")
storeRoutes.Use(middleware.JWTAuth(cfg, db))
storeRoutes.Use(middleware.RequireAdminRoles())

// 員工管理相關 - 經理級別以上
staffRoutes := router.Group("/api/v1/staff")
staffRoutes.Use(middleware.JWTAuth(cfg, db))
staffRoutes.Use(middleware.RequireManagerOrAbove())

// 預約管理相關 - 經理和美甲師
bookingRoutes := router.Group("/api/v1/booking")
bookingRoutes.Use(middleware.JWTAuth(cfg, db))
bookingRoutes.Use(middleware.RequireRoles(staff.RoleManager, staff.RoleStylist))
```

## 技術細節

### 權限驗證流程

1. 檢查 `Gin context` 中是否存在員工上下文（由 `JWT middleware` 設置）
2. 如果不存在，返回 401 認證失敗
3. 檢查員工角色是否在允許的角色列表中
4. 如果不在列表中，返回 403 權限不足
5. 如果驗證通過，繼續執行下一個 middleware 或處理器

### 擴展性考慮

- 支援多角色組合，靈活配置權限
- 便利函數簡化常見權限場景
- 易於添加新的權限檢查邏輯
- 與現有 JWT 認證系統完全兼容
