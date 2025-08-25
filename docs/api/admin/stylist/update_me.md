## User Story

作為一位員工，我希望可以更新自己的美甲師個人資料（stylists），讓顧客可以看到我最新的專長與風格。

---

## Endpoint

**PATCH** `/api/admin/stylists/me`

---

## 說明

- 僅允許員工更新自己的美甲師資料。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "stylistName": "Jane 美甲師",
  "goodAtShapes": ["方形", "橢圓形"],
  "goodAtColors": ["白色系"],
  "goodAtStyles": ["簡約", "法式"],
  "isIntrovert": true
}
```

### 驗證規則

| 欄位         | 必填 | 其他規則                                                                                           | 說明     |
| ------------ | ---- | -------------------------------------------------------------------------------------------------- | -------- |
| stylistName  | 否   | <li>不能為空字串<li>最大長度50字元                                                                 | 顯示名稱 |
| goodAtShapes | 否   | <li>最多20項<li>值只能為 方形 方圓形 橢圓形 圓形 圓尖形 尖形 梯形                                  | 擅長指型 |
| goodAtColors | 否   | <li>最多20項<li>值只能為 白色系 裸色系 粉色系 紅色系 橘色系 大地色系 綠色系 藍色系 紫色系 黑色系   | 擅長色系 |
| goodAtStyles | 否   | <li>最多20項<li>值只能為 暈染 手繪 貓眼 鏡面 可愛 法式 漸層 氣質溫柔 個性 日系 簡約 優雅 典雅 小眾 | 擅長款式 |
| isIntrovert  | 否   |                                                                                                    | 是否I人  |

- 至少需要提供一個欄位進行更新

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "18000000001",
    "staffUserId": "13984392823",
    "name": "Jane 美甲師",
    "goodAtShapes": ["方形", "橢圓形"],
    "goodAtColors": ["粉嫩系"],
    "goodAtStyles": ["簡約", "法式"],
    "isIntrovert": true,
    "createdAt": "2025-06-01T08:00:00+08:00",
    "updatedAt": "2025-06-01T08:00:00+08:00"
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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                                  |
| ------ | ------ | ----------------------- | ------------------------------------- |
| 401    | E1002  | AuthTokenInvalid        | 無效的 accessToken，請重新登入        |
| 401    | E1003  | AuthTokenMissing        | accessToken 缺失，請重新登入          |
| 401    | E1004  | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入      |
| 401    | E1005  | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006  | AuthContextMissing      | 未找到使用者認證資訊，請重新登入      |
| 403    | E3001  | AuthPermissionDenied    | 權限不足，無法執行此操作              |
| 400    | E2001  | ValJsonFormat           | JSON 格式錯誤，請檢查                 |
| 400    | E2002  | AuthPathParamMissing    | 路徑參數缺失，請檢查                  |
| 400    | E2004  | AuthParamTypeConversion | 參數類型轉換失敗                      |
| 400    | E2020  | ValFieldRequired        | {field} 為必填項目                    |
| 400    | E2023  | ValFieldMinNumber       | {field} 最小值為 {param}              |
| 400    | E2024  | ValFieldMaxLength       | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2026  | ValFieldMaxNumber       | {field} 最大值為 {param}              |
| 400    | E2036  | ValFieldNoBlank         | {field} 不能為空字串                  |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                        |

---

## 資料表

- `stylists`
- `staff_users`

---

## Service 邏輯

1. 檢查 `stylists` 資料是否存在。
2. 更新 `stylists` 的指定欄位。
3. 回傳更新後的 `stylist` 資料。

---

## 注意事項

- 會回傳更新後的資料，不需要在呼叫一次 `GET /api/admin/staff/me` 來取得資料。
