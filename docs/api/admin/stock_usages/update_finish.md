## User Story

作為一位管理員，我希望能更新產品庫存使用紀錄，方便維護庫存使用紀錄。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}/stock-usages/{stockUsageId}/finish`

---

## 說明

- 提供後台管理員更新產品庫存使用紀錄為完成狀態功能。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` `MANAGER` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數         | 說明           |
| ------------ | -------------- |
| stockUsageId | 庫存使用紀錄ID |

### Body 範例

```json
{
  "usageEndedAt": "2025-01-01"
}
```

### 驗證規則

| 欄位         | 必填 | 其他規則             | 說明         |
| ------------ | ---- | -------------------- | ------------ |
| usageEndedAt | 是   | <li>格式是YYYY-MM-DD | 使用結束日期 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "9000000001"
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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                                                |
| ------ | -------- | ----------------------- | --------------------------------------------------- |
| 401    | E1002    | AuthTokenInvalid        | 無效的 accessToken，請重新登入                      |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入                        |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入                    |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入                    |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入                    |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作                            |
| 400    | E2001    | ValJsonFormat           | JSON 格式錯誤，請檢查                               |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查                                |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                                    |
| 400    | E2020    | ValFieldRequired        | {field} 為必填項目                                  |
| 400    | E2033    | ValFieldDateFormat      | {field} 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD) |
| 400    | E3STU002 | StockUsageNotInUse      | 庫存使用紀錄已完成，無法再次完成                    |
| 404    | E3STU001 | StockUsageNotFound      | 庫存使用紀錄不存在或已被刪除                        |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試                            |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                                      |

---

## 資料表

- `stock_usages`

---

## Service 邏輯

1. 確認門市存取權限。
2. 確認庫存使用紀錄是否存在。
4. 確認庫存使用紀錄是否為正在使用中。
5. 更新 `stock_usages` 資料。
6. 回傳更新結果。
