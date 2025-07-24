## User Story

作為客戶，我想要取得我自己的資料，方便於前台顯示或個人中心查閱。

---

## Endpoint

**GET** `/api/customers/me`

---

## 說明

- 僅支援已登入客戶（access token 驗證）。
- 回傳當前登入客戶的完整資料。
- 適用於「我的帳號」、「個人中心」等場景。

---

## 權限

- 僅客戶本人可查詢（JWT 驗證）。

---

## Request

### Header

```http
Authorization: Bearer <access_token>
```

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "1000000001",
    "name": "王小美",
    "phone": "0912345678",
    "birthday": "1992-02-29",
    "city": "台北市",
    "favoriteShapes": ["方形"],
    "favoriteColors": ["粉色"],
    "favoriteStyles": ["法式"],
    "isIntrovert": false,
    "referralSource": ["朋友推薦"],
    "referrer": "黃小姐",
    "customerNote": "容易指緣乾裂"
  }
}
```

### 失敗

#### 401 Unauthorized - 未登入/Token 失效

```json
{
  "message": "無效的 accessToken"
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

---

## Service 邏輯

1. 查詢並回傳該客戶的完整資料。

---

## 注意事項

- 僅允許本人查詢。

