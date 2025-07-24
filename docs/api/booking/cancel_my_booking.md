## User Story

作為顧客，我希望可以取消自己的預約，並傳入取消原因，讓店家知悉我的狀況。

---

## Endpoint

**PATCH** `/api/bookings/{bookingId}/cancel`

---

## 說明

- 僅支援已登入顧客（access token 驗證）。
- 僅允許本人取消自己的預約。
- 可傳入取消原因（文字），供後台記錄與分析。
- 預約狀態將變更為 CANCELLED。

---

## 權限

- 僅顧客本人可操作（JWT 驗證）。

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Path Parameter

| 參數      | 說明   |
| --------- | ------ |
| bookingId | 預約ID |

### Body

```json
{
  "cancelReason": "臨時有事無法前往，抱歉！"
}
```

### 驗證規則

| 欄位         | 規則                  | 說明     |
| ------------ | --------------------- | -------- |
| cancelReason | <li>選填<li>最長100字 | 取消原因 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "5000000001",
    "status": "CANCELLED",
    "cancelReason": "臨時有事無法前往，抱歉！"
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗"
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 權限不足

```json
{
  "message": "權限不足，僅限本人操作"
}
```

#### 404 Not Found - 預約不存在

```json
{
  "message": "預約不存在或已被刪除"
}
```

#### 409 Conflict - 已取消或不可取消狀態

```json
{
  "message": "預約已取消或不可再取消"
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

- `bookings`

---

## Service 邏輯

1. 驗證預約是否存在且屬於本人，且狀態為 `SCHEDULED`。
2. 記錄取消原因，變更狀態為 `CANCELLED`。
3. 回傳結果。

---

## 注意事項

- 僅支援本人預約取消。
- 取消時若有傳入原因，則記錄取消原因。
- 狀態不可重複取消。

