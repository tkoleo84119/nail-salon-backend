## User Story

作為一位員工，我希望能夠登出系統。

---

## Endpoint

**POST** `/api/admin/auth/token/revoke`

---

## 說明

- 提供後台員工登出功能。
- 登出時無論傳入的 `refreshToken` 是否正確，皆會回傳 `success`，確保前端登出體驗簡潔且安全。
- 此為後台管理端使用，與前台顧客端獨立。

---

## 權限

- 不須預先認證。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "refreshToken": "your_refresh_token_here"
}
```

### 驗證規則

| 欄位         | 必填 | 其他規則        | 說明               |
| ------------ | ---- | --------------- | ------------------ |
| refreshToken | 是   | 最大長度500字元 | 員工 Refresh Token |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "success": true
  }
}
```

---

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

| 狀態碼 | 錯誤碼 | 常數名稱          | 說明                                  |
| ------ | ------ | ----------------- | ------------------------------------- |
| 400    | E2001  | ValJsonFormat     | JSON 格式錯誤，請檢查                 |
| 400    | E2020  | ValFieldRequired  | {field} 為必填項目                    |
| 400    | E2024  | ValFieldMaxLength | {field} 長度最多只能有 {param} 個字元 |
| 500    | E9001  | SysInternalError  | 系統發生錯誤，請稍後再試              |

---

## 實作與流程

### Service 邏輯

1. 嘗試將資料庫對應 refreshToken 設為 `is_revoked = true`。
2. 不論資料庫操作結果，皆回傳 `{ "success": true }`。

### 資料表

- `staff_user_tokens`

---

## 注意事項

- 登出行為不會影響 `Access Token`。
- 無論 `refreshToken` 是否有效或已被撤銷，皆回傳成功。
