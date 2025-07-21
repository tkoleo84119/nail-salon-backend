## User Story

作為一個客戶，我希望可以編輯自己的客戶資料，方便保持資訊正確。

---

## Endpoint

**PATCH** `/api/customers/me`

---

## 說明

- 僅支援已登入客戶（access token 驗證）。
- 客戶可編輯自己的基本資料：姓名、電話、生日、常用偏好、備註。
- 皆為選填。

---

## 權限

- 僅客戶本人可編輯自己的資料（JWT 驗證）。

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Body

```json
{
  "name": "王小美",
  "phone": "0912345678",
  "birthday": "1992-02-29",
  "city": "台北市",
  "favoriteShapes": ["方形"],
  "favoriteColors": ["粉色"],
  "favoriteStyles": ["法式"],
  "isIntrovert": true,
  "customerNote": "容易指緣乾裂"
}
```

- 欄位皆為選填，但至少需有一項。

### 驗證規則

| 欄位           | 規則                                                        | 說明       |
| -------------- | ----------------------------------------------------------- | ---------- |
| name           | <li>選填<li>長度最小是1<li>長度最大是100                    | 姓名       |
| phone          | <li>選填<li>長度最小是1<li>長度最大是20<li>格式是09xxxxxxxx | 電話       |
| birthday       | <li>選填<li>格式是yyyy-MM-dd                                | 生日       |
| city           | <li>選填<li>長度最小是1<li>長度最大是100                    | 城市       |
| favoriteShapes | <li>選填<li>陣列                                            | 喜歡的指形 |
| favoriteColors | <li>選填<li>陣列                                            | 喜歡色系   |
| favoriteStyles | <li>選填<li>陣列                                            | 喜歡款式   |
| isIntrovert    | <li>選填<li>布林值                                          | 是否是I人  |
| customerNote   | <li>選填<li>長度最小是1<li>長度最大是1000                   | 個人備註   |

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
    "isIntrovert": true,
    "referralSource": ["朋友介紹", "網路廣告"],
    "referrer": "1000000001",
    "customerNote": "容易指緣乾裂"
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "phone": "必須為有效的台灣手機號碼格式 (例: 09xxxxxxxx)"
  }
}
```

#### 401 Unauthorized - 未登入/Token 失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 權限不足

```json
{
  "message": "權限不足，無法執行此操作"
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

1. 若有傳入 birthday，則驗證格式是否為 yyyy-MM-dd。
2. 確認 customer 資料是否存在。
3. 更新 customer 資料。
4. 回傳更新後資料。

---

## 注意事項

- 僅允許本人編輯。
- 一定要有至少一項欄位被更新。

