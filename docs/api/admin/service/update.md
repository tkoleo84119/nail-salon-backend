## User Story

作為一位管理員，我希望能更新服務項目，方便即時維護可預約之美甲服務內容。

---

## Endpoint

**PATCH** `/api/admin/services/{serviceId}`

---

## 說明

- 可更新名稱、價格、操作時間、是否為附加服務、顯示狀態、啟用狀態、備註。
- 服務名稱須唯一(不包含自己)。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數      | 說明       |
| --------- | ---------- |
| serviceId | 服務項目ID |

### Body 範例

```json
{
  "sortOrder": 1,
  "name": "凝膠足部單色",
  "price": 1400,
  "durationMinutes": 75,
  "isAddon": false,
  "isVisible": true,
  "isActive": true,
  "note": "足部基礎保養"
}
```

### 驗證規則

### 驗證規則

| 欄位            | 必填 | 其他規則                            | 說明     |
| --------------- | ---- | ----------------------------------- | -------- |
| sortOrder       | 否   | <li>最小值0<li>最大值1000000        | 排序序號 |
| name            | 否   | <li>不能為空字串<li>最大長度100字元 | 服務名稱 |
| price           | 否   | <li>最小值0<li>最大值1000000        | 價格     |
| durationMinutes | 否   | <li>最小值0<li>最大值1440           | 操作分鐘 |
| isAddon         | 否   |                                     | 附加服務 |
| isVisible       | 否   |                                     | 可見狀態 |
| isActive        | 否   |                                     | 啟用狀態 |
| note            | 否   | <li>最大長度255                     | 備註     |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "9000000001",
    "sortOrder": 1,
    "name": "凝膠足部單色",
    "price": 1400,
    "durationMinutes": 75,
    "isAddon": false,
    "isVisible": true,
    "isActive": true,
    "note": "足部基礎保養",
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
| 401    | E1002    | AuthTokenInvalid        | 無效的 accessToken，請重新登入        |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入          |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入      |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入      |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作              |
| 400    | E2001    | ValJsonFormat           | JSON 格式錯誤，請檢查                 |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查                  |
| 400    | E2003    | ValAllFieldsEmpty       | 至少需要提供一個欄位進行更新          |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                      |
| 400    | E2023    | ValFieldMinNumber       | {field} 最小值為 {param}              |
| 400    | E2024    | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2026    | ValFieldMaxNumber       | {field} 最大值為 {param}              |
| 400    | E2036    | ValFieldNoBlank         | {field} 不能為空字串                  |
| 404    | E3SER004 | ServiceNotFound         | 服務不存在或已被刪除                  |
| 409    | E3SER005 | ServiceAlreadyExists    | 服務已存在                            |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                        |

---

## 資料表

- `services`

---

## Service 邏輯

1. 驗證 `serviceId` 是否存在。
2. 若有更新 `name`，則驗證名稱是否唯一（不包含自己）。
3. 更新 `services` 資料。
4. 回傳更新結果。

---

## 注意事項

- 服務名稱不可重複。
- 設定為不可見或未啟用時，前台不可被預約。
