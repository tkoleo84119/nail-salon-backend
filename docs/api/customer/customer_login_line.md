## User Story

作為一個客戶，當我以 LINE 登入並將 idToken 傳送給後端時，
系統會自動判斷我是否已註冊本應用程式，
若已註冊則回傳 `access token` 和 `refresh token`，否則引導我完成註冊流程（回傳 `needRegister` 為 `true`，非錯誤狀態碼）。

---

## Endpoint

**POST** `/api/auth/customer/line/login`

---

## 說明

- 用戶端取得 LINE idToken 後，傳給本 API 進行登入/註冊判斷。
- 若 LINE idToken 驗證失敗，回傳 400。
- 若已註冊本系統，直接發 access token / refresh token。
- 若尚未註冊，回傳需註冊訊息與 LINE 基本資訊，前端據此引導註冊流程。

---

## Request

### Header

```http
Content-Type: application/json
```

### Body

```json
{
  "idToken": "eyJraWQiOiJ..."
}
```

| 欄位    | 規則                                     | 說明         |
| ------- | ---------------------------------------- | ------------ |
| idToken | <li>必填<li>長度最小是1<li>長度最大是500 | LINE idToken |

---

## Response

### 已註冊：200 OK

```json
{
  "data": {
    "needRegister": false,
    "accessToken": "1234567890",
    "refreshToken": "1234567890",
    "customer": {
      "id": "1000000001",
      "name": "小美",
      "phone": "09xxxxxxxx"
    }
  }
}
```

### 尚未註冊：200 OK

```json
{
  "needRegister": true,
  "lineProfile": {
    "providerUid": "U12345678",
    "name": "Mei",
    "email": "mei@example.com" // 可能為空
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤 / idToken 無效

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "idToken": "idToken為必填項目"
  }
}
```

```json
{
  "message": "idToken 驗證失敗"
}
```

#### 500 Internal Server Error

```json
{
  "message": "系統發生錯誤，請稍後再試"
}
```

---

## 資料表

- `customers`
- `customer_auths`
- `customer_tokens`

---

## Service 邏輯

1. 驗證 idToken 合法性，解析 providerUid。
2. 查詢 providerUid 是否已綁定客戶。
   - 已註冊：產生 access token / refresh token。
   - 未註冊：回傳需註冊及 LINE profile。
3. 回傳對應資訊。

---

## 注意事項

- 登入與註冊流程需區分處理，請前端依據 `needRegister` 做引導。
