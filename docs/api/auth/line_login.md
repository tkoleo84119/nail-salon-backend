## User Story

作為一個客戶，當我以 LINE 登入並將 idToken 傳送給後端時，
系統會自動判斷我是否已註冊本應用程式，
若已註冊則回傳 `accessToken` 和 `refreshToken`，否則引導我完成註冊流程（回傳 `needRegister` 為 `true`，非錯誤狀態碼）。

---

## Endpoint

**POST** `/api/auth/line/login`

---

## 說明

- 用戶端取得 LINE idToken 後，傳給本 API 進行登入/註冊判斷。
- 若已註冊本系統，直接發 `accessToken` 和 `refreshToken`。
- 若尚未註冊，回傳需註冊訊息與 LINE 基本資訊，前端據此引導註冊流程。

---

## 權限

- 不須預先認證

---

## Request

### Header

- Content-Type: application/json

### Body 範例

```json
{
  "idToken": "eyJraWQiOiJ..."
}
```

### 驗證規則

| 欄位    | 必填 | 其他規則                             | 說明         |
| ------- | ---- | ------------------------------------ | ------------ |
| idToken | 是   | <li>不能為空字串<li>最大長度2000字元 | LINE idToken |

---

## Response

### 已註冊：200 OK

```json
{
  "data": {
    "needRegister": false,
    "needCheckTerms": false, // 是否需要使用者確認條款
    "accessToken": "1234567890",
    "refreshToken": "1234567890",
    "expiresIn": 3600
  }
}
```

### 尚未註冊：200 OK

```json
{
  "data": {
    "needRegister": true,
    "lineProfile": {
      "providerUid": "U12345678",
      "name": "Mei",
      "email": "mei@example.com" // 可能為空
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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                                  |
| ------ | ------ | ----------------------- | ------------------------------------- |
| 401    | E1007  | AuthLineTokenInvalid    | Line idToken 驗證失敗，請重新登入     |
| 401    | E1008  | AuthLineTokenExpired    | Line idToken 已過期，請重新登入       |
| 400    | E2001  | ValJsonFormat           | JSON 格式錯誤，請檢查                 |
| 400    | E2004  | ValTypeConversionFailed | 參數類型轉換失敗                      |
| 400    | E2020  | ValFieldRequired        | {field} 為必填項目                    |
| 400    | E2036  | ValFieldNoBlank         | {field} 不能為空字串                  |
| 400    | E2024  | ValFieldMaxLength       | {field} 長度最多只能有 {param} 個字元 |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                        |

---

## 資料表

- `customers`
- `customer_tokens`

---

## Service 邏輯

1. 呼叫 LINE 驗證 `idToken` 合法性，取得 `providerUid`。
2. 根據 `providerUid` 查詢 `customers` 資料。
   - 未註冊：回傳需註冊及 LINE profile 及 `needRegister` 為 `true`。
3. 產生發 `access token`、`refresh token`。
4. 檢查客戶是否有更新 `line_name`，若有則更新 (避免 LINE 名稱更新但資料庫未更新)。
5. 回傳 `access token`、`refresh token`。

---

## 注意事項

- 登入與註冊流程需區分處理，前端依據 `needRegister` 做引導。
