## User Story

作為一位管理員（`SUPER_ADMIN` / `ADMIN`），我希望能更新服務項目，方便即時維護可預約之美甲服務內容。

---

## Endpoint

**PATCH** `/api/services/{serviceId}`

---

## 說明

- 僅限 `SUPER_ADMIN`、`ADMIN` 可更新服務項目。
- 可更新名稱、價格、操作時間、是否為附加服務、顯示狀態、啟用狀態、備註。
- 服務名稱須唯一(不包含自己)。
- 欄位皆為選填，但至少需有一項。

---

## 權限

- 僅 `SUPER_ADMIN`、`ADMIN` 可操作。

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Path Parameter

| 參數      | 說明       |
| --------- | ---------- |
| serviceId | 服務項目ID |

### Body

```json
{
  "name": "凝膠足部單色",
  "price": 1400,
  "durationMinutes": 75,
  "isAddon": false,
  "isVisible": true,
  "isActive": true,
  "note": "足部基礎保養"
}
```

### 驗證規則

| 欄位            | 規則                                         | 說明     |
| --------------- | -------------------------------------------- | -------- |
| name            | <li>選填<li>長度大於1<li>長度小於100<li>唯一 | 服務名稱 |
| price           | <li>選填<li>數字最小是0                      | 價格     |
| durationMinutes | <li>選填<li>數字最小是0<li>小於1440          | 操作分鐘 |
| isAddon         | <li>選填<li>布林值                           | 附加服務 |
| isVisible       | <li>選填<li>布林值                           | 可見狀態 |
| isActive        | <li>選填<li>布林值                           | 啟用狀態 |
| note            | <li>選填<li>長度小於255                      | 備註     |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "9000000001",
    "name": "凝膠足部單色",
    "price": 1400,
    "durationMinutes": 75,
    "isAddon": false,
    "isVisible": true,
    "isActive": true,
    "note": "足部基礎保養"
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "name": "name最小長度為1"
  }
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
  "message": "權限不足，僅限管理員操作"
}
```

#### 404 Not Found - 服務不存在

```json
{
  "message": "服務不存在或已被刪除"
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

- `services`

---

## Service 邏輯

1. 驗證 `serviceId` 是否存在。
2. 若有更新 name，則驗證名稱是否唯一（不包含自己）。
3. 更新 `services` 資料。
4. 回傳更新結果。

---

## 注意事項

- 服務名稱不可重複。
- 支援「主服務」及「附加服務」類型。
- 設定為不可見或未啟用時，前台不可被預約。

