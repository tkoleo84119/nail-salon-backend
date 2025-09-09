## User Story

作為顧客，我希望可以確認條款，方便我使用本應用程式。

---

## Endpoint

**POST** `/api/auth/accept-term`

---

## 說明

- 提供顧客確認條款。

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
  "termsVersion": "v1"
}
```

### 驗證規則

| 欄位         | 必填 | 其他規則        | 說明     |
| ------------ | ---- | --------------- | -------- |
| termsVersion | 是   | <li>值只能為 v1 | 條款版本 |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "5000000001"
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
| 401    | E1006  | AuthContextMissing   | 未找到使用者認證資訊，請重新登入 |
| 401    | E1011  | AuthCustomerFailed   | 未找到有效的顧客資訊，請重新登入 |
| 400    | E2001  | ValJSONFormatError   | JSON 格式錯誤，請檢查            |
| 400    | E2020  | ValFieldRequired     | {field} 為必填項目               |
| 500    | E9001  | SysInternalError     | 系統發生錯誤，請稍後再試         |
| 500    | E9002  | SysDatabaseError     | 資料庫操作失敗                   |

---

## 資料表

- `customer_terms_acceptance`

---

## Service 邏輯

1. 驗證是否已經確認過條款。
2. 若沒有，則新增條款確認資料（`customer_terms_acceptance`）。
3. 若有，則不新增。
4. 回傳條款確認資訊。

---

## 注意事項

- 僅支援本人確認條款。
