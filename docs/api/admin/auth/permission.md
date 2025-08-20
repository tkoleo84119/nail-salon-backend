## User Story

作為一位員工，我希望能夠取得自己的權限相關資料，例如：可存取的店家、角色。

---

## Endpoint

**GET** `/api/admin/auth/permission`

---

## 說明

- 員工提供 `accessToken`，若有效則回傳自己的權限相關資料。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "139842394",
    "username": "admin001",
    "role": "ADMIN",
    "storeAccess": [
      {
        "id": "1",
        "name": "門市1"
      }
    ]
  }
}
```

### 錯誤處理

全部 API 皆回傳如下結構，請參考錯誤總覽。

```json
{
  "errors": [
    {
      "code": "EXXXX",
      "message": "錯誤訊息",
      "field": "錯誤欄位名稱"
    }
  ]
}
```

- 欄位說明：
  - errors: 錯誤陣列（支援多筆同時回報）
  - code: 錯誤代碼，唯一對應每種錯誤
  - message: 中文錯誤訊息（可參照錯誤總覽）
  - field: 參數欄位名稱（僅部分驗證錯誤有）

| 狀態碼 | 錯誤碼 | 常數名稱             | 說明                             |
| ------ | ------ | -------------------- | -------------------------------- |
| 401    | E1002  | AuthTokenInvalid     | 無效的 accessToken，請重新登入   |
| 401    | E1003  | AuthTokenMissing     | accessToken 缺失，請重新登入     |
| 401    | E1004  | AuthTokenFormatError | accessToken 格式錯誤，請重新登入 |
| 401    | E1005  | AuthStaffFailed      | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006  | AuthContextMissing   | 未找到使用者認證資訊，請重新登入 |
| 500    | E9001  | SysInternalError     | 系統發生錯誤，請稍後再試         |
| 500    | E9002  | SysDatabaseError     | 資料庫操作失敗                   |

---

## 實作與流程

### Service 邏輯

1. 從 `StaffContext` 取得員工資料
2. 回傳員工資料與可存取的店家

---

## 注意事項

- 直接從 `StaffContext` 取得員工資料，不會再從資料庫取得，因為 `StaffContext` 已經驗證過了。
