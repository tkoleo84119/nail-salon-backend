## User Story

作為顧客，我希望能夠取得特定門市資料，方便查詢分店資訊。

---

## Endpoint

**GET** `/api/stores/{storeId}`

---

## 說明

- 支援已登入顧客查詢單一門市資訊。
- 僅回傳 `is_active=true` 的門市。
- 適用於門市介紹等場景。

---

## 權限

- 僅顧客可查詢（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明   |
| ------- | ------ |
| storeId | 門市ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "8000000001",
    "name": "大安旗艦店",
    "address": "台北市大安區復興南路一段100號",
    "phone": "02-1234-5678"
  }
}
```

### 失敗

#### 404 Not Found - 門市不存在/未啟用

```json
{
  "message": "查無此門市或門市未啟用"
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

- `stores`

---

## Service 邏輯

1. 查詢 storeId 對應之 `is_active=true` 門市。
2. 回傳門市資訊。

---

## 注意事項

- 僅回傳啟用門市。

