## User Story

作為一個客戶，當我要完成 LINE 註冊時，系統會將我的資訊存入後，回傳 access token、refresh token 等。

---

## Endpoint

**POST** `/api/auth/customer/line/register`

---

## 說明

- 用戶在 LINE 登入後，若尚未註冊，需呼叫本 API 完成註冊流程。
- 完成註冊後，自動發 access token 與 refresh token。
- 需帶入從 LINE profile 取得的必要資訊與額外註冊欄位。

---

## Request

### Header

```http
Content-Type: application/json
```

### Body

```json
{
  "idToken": "eyJraWQiOiJ...",
  "name": "小美",
  "phone": "09xxxxxxxx",
  "birthday": "1990-01-01",
  "city": "台北市",
  "favorite_shapes": ["圓形", "方形"],
  "favorite_colors": ["黑色", "白色"],
  "favorite_styles": ["自然", "韓式"],
  "is_introvert": true,
  "referral_source": ["朋友介紹", "網路廣告"],
  "referrer": "1000000001",
  "customer_note": "這是客戶的備註",
}
```

| 欄位            | 規則                                                        | 說明             |
| --------------- | ----------------------------------------------------------- | ---------------- |
| idToken         | <li>必填<li>長度最小是1<li>長度最大是500                    | LINE idToken     |
| name            | <li>必填<li>長度最小是1<li>長度最大是100                    | 姓名             |
| phone           | <li>必填<li>長度最小是1<li>長度最大是20<li>格式是09xxxxxxxx | 電話             |
| birthday        | <li>必填<li>格式是yyyy-MM-dd                                | 生日             |
| city            | <li>選填<li>長度最小是1<li>長度最大是100                    | 城市             |
| favorite_shapes | <li>選填<li>陣列<li>長度最小是1                             | 喜歡的指形       |
| favorite_colors | <li>選填<li>陣列<li>長度最小是1                             | 喜歡的色系       |
| favorite_styles | <li>選填<li>陣列<li>長度最小是1                             | 喜歡的款式       |
| is_introvert    | <li>選填<li>布林值                                          | 是否是I人        |
| referral_source | <li>選填<li>陣列<li>長度最小是1                             | 推薦來源         |
| referrer        | <li>選填<li>長度最小是1                                     | 推薦人           |
| customer_note   | <li>選填<li>長度最小是1<li>長度最大是1000                   | 使用者自己的備註 |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "accessToken": "...",
    "refreshToken": "...",
    "customer": {
      "id": "1000000001",
      "name": "小美",
      "phone": "09xxxxxxxx",
      "birthday": "1990-01-01"
    }
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤 / idToken 無效 / 欄位未填

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "name": "name為必填項目"
  }
}
```

#### 409 Conflict - LINE 帳號已註冊

```json
{
  "message": "該 LINE 帳號已註冊過，請直接登入"
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
2. 驗證該 LINE UID 是否已註冊（重複則 409）。
3. 建立 customer 資料與對應 customer_auths。
4. 發 access token、refresh token。
5. 回傳客戶資訊。

---

## 注意事項

- 電話格式必須為 09xxxxxxxx。
- 生日格式必須為 yyyy-MM-dd。
- accessToken / refreshToken 請依安全原則管理。

---
