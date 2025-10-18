## User Story

作為員工，我希望可以查詢特定服務資料，以便管理或顯示詳細內容。

---

## Endpoint

**GET** `/api/admin/services/{serviceId}`

---

## 說明

- 用於查詢特定服務的詳細資訊。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameters

| 參數      | 說明    |
| --------- | ------- |
| serviceId | 服務 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "9000000001",
    "sortOrder": 1,
    "name": "手部單色",
    "durationMinutes": 60,
    "price": 1200,
    "isAddon": false,
    "isActive": true,
    "isVisible": true,
    "note": "含修型保養",
    "createdAt": "2025-01-01T00:00:00+08:00",
    "updatedAt": "2025-01-01T00:00:00+08:00"
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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                             |
| ------ | -------- | ----------------------- | -------------------------------- |
| 401    | E1002  | AuthTokenInvalid       | 無效的 accessToken，請重新登入   |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入     |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作         |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查             |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                 |
| 404    | E3SER004 | ServiceNotFound         | 服務不存在或已被刪除             |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                   |

---

## 資料表

- `services`

---

## Service 邏輯

1. 查詢 `services` 表中該筆服務是否存在。
2. 不存在則回傳 `404 Not Found`。
3. 存在則回傳該筆服務詳細內容。
