## User Story

作為一位管理員，我希望能新增服務項目，方便維護可預約之美甲服務。

---

## Endpoint

**POST** `/api/admin/services`

---

## 說明

- 提供後台管理員新增服務項目功能。
- 服務名稱須唯一。
- 可設定價格、操作時間、是否為附加服務、顯示狀態與備註。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "name": "凝膠手部單色",
  "price": 1200,
  "durationMinutes": 60,
  "isAddon": false,
  "isVisible": true,
  "note": "含基礎修型保養"
}
```

### 驗證規則

| 欄位            | 必填 | 其他規則                     | 說明     |
| --------------- | ---- | ---------------------------- | -------- |
| name            | 是   | <li>最大長度100字元          | 服務名稱 |
| price           | 是   | <li>最小值0<li>最大值1000000 | 價格     |
| durationMinutes | 是   | <li>最小值0<li>最大值1440    | 操作分鐘 |
| isAddon         | 是   | <li>布林值                   | 附加服務 |
| isVisible       | 是   | <li>布林值                   | 可見狀態 |
| note            | 選填 | <li>最大長度255              | 備註     |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "9000000001",
    "name": "凝膠手部單色",
    "price": 1200,
    "durationMinutes": 60,
    "isAddon": false,
    "isVisible": true,
    "isActive": true,
    "note": "含基礎修型保養",
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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                                  |
| ------ | -------- | ----------------------- | ------------------------------------- |
| 401    | E1002    | AuthInvalidCredentials  | 無效的 accessToken，請重新登入        |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入          |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入      |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入      |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作              |
| 400    | E2001    | ValJsonFormat           | JSON 格式錯誤，請檢查                 |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查                  |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                      |
| 400    | E2020    | ValFieldRequired        | {field} 為必填項目                    |
| 400    | E2023    | ValFieldMinNumber       | {field} 最小值為 {param}              |
| 400    | E2024    | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2026    | ValFieldMaxNumber       | {field} 最大值為 {param}              |
| 400    | E2029    | ValFieldBoolean         | {field} 必須是布林值                  |
| 409    | E3SER005 | ServiceAlreadyExists    | 服務已存在                            |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                        |

---

## 資料表

- `services`

---

## Service 邏輯

1. 驗證角色是否為 `SUPER_ADMIN` 或 `ADMIN`。
2. 驗證 `name` 是否唯一。
3. 建立 `services` 資料。
4. 回傳新增結果。

---

## 注意事項

- 服務名稱不可重複。
- 設定為不可見或未啟用時，客戶不可從前台預約。
