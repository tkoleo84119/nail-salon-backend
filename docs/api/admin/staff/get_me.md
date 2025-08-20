## User Story

作為員工，我希望可以取得自己的員工資料，方便查看登入帳號與基本資訊。

---

## Endpoint

**GET** `/api/admin/staff/me`

---

## 說明

- 提供當前登入員工自己的基本資料。
- 若該員工同時為美甲師（`stylists.staff_user_id`），則一併回傳 stylist 資訊。
- 根據 JWT 中的 `staff_user_id` 查詢對應資料。

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
    "id": "6000000002",
    "username": "stylist88",
    "email": "s88@salon.com",
    "role": "STYLIST",
    "isActive": true,
    "createdAt": "2025-06-01T08:00:00+08:00",
    "updatedAt": "2025-06-01T08:00:00+08:00",
    "stylist": {
      "id": "7000000001",
      "name": "Bella",
      "goodAtShapes": ["方形"],
      "goodAtColors": ["粉色系"],
      "goodAtStyles": ["簡約風"],
      "isIntrovert": false,
      "createdAt": "2025-06-01T08:00:00+08:00",
      "updatedAt": "2025-06-01T08:00:00+08:00"
    }
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
| 404    | E3STA005 | StaffNotFound           | 員工帳號不存在                   |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                   |

---

## 資料表

- `staff_users`
- `stylists`

---

## Service 邏輯

1. 查詢 `staff_users` 表是否存在該使用者。
2. 若存在，則回傳基本資訊（帳號、信箱、角色、狀態、建立時間）。
3. 若該員工同時為美甲師（`stylists.staff_user_id`），則一併回傳 stylist 資訊。

---

## 注意事項

- 僅可查詢自己的帳號，無其他查詢條件。
