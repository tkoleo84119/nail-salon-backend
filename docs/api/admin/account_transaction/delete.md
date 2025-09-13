## User Story

作為一位管理員，我希望能刪除最新一筆帳戶紀錄，方便即時維護帳戶紀錄。

---

## Endpoint

**DELETE** `/api/admin/stores/{storeId}/accounts/{accountId}/transactions/latest`

---

## 說明

- 可刪除最新一筆帳戶紀錄。

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

| 參數      | 說明   |
| --------- | ------ |
| storeId   | 門市ID |
| accountId | 帳戶ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "deleted": "6000000011"
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

| 狀態碼 | 錯誤碼  | 常數名稱                   | 說明                             |
| ------ | ------- | -------------------------- | -------------------------------- |
| 401    | E1002   | AuthTokenInvalid           | 無效的 accessToken，請重新登入   |
| 401    | E1003   | AuthTokenMissing           | accessToken 缺失，請重新登入     |
| 401    | E1004   | AuthTokenFormatError       | accessToken 格式錯誤，請重新登入 |
| 401    | E1005   | AuthStaffFailed            | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006   | AuthContextMissing         | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010   | AuthPermissionDenied       | 權限不足，無法執行此操作         |
| 400    | E2001   | ValJsonFormat              | JSON 格式錯誤，請檢查            |
| 400    | E2002   | ValPathParamMissing        | 路徑參數缺失，請檢查             |
| 400    | E2004   | ValTypeConversionFailed    | 參數類型轉換失敗                 |
| 400    | E3ACC02 | AccountNotBelongToStore    | 帳戶不屬於指定的門市             |
| 404    | E3ACC01 | AccountNotFound            | 帳戶不存在或已被刪除             |
| 404    | E3ACC03 | AccountTransactionNotFound | 帳戶交易紀錄不存在或已被刪除     |
| 500    | E9001   | SysInternalError           | 系統發生錯誤，請稍後再試         |
| 500    | E9002   | SysDatabaseError           | 資料庫操作失敗                   |

---

## 資料表

- `accounts`
- `account_transactions`

---

## Service 邏輯

1. 驗證 `account` 是否存在。
2. 刪除最新一筆 `account_transactions` 資料。
3. 回傳刪除結果。
