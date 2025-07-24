## User Story

作為一位管理員（`SUPER_ADMIN` / `ADMIN`），我希望能新增服務項目，方便維護可預約之美甲服務。

---

## Endpoint

**POST** `/api/admin/services`

---

## 說明

- 僅限 `SUPER_ADMIN`、`ADMIN` 可建立新服務。
- 服務名稱須唯一。
- 可設定價格、操作時間、是否為附加服務、顯示狀態與備註。

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

### Body

```json
{
  "name": "凝膠手部單色",
  "price": 1200,
  "durationMinutes": 60,
  "isAddon": false,
  "isVisible": true,
  "note": "含基礎修型保養"
}
```

### 驗證規則

| 欄位            | 規則                                         | 說明     |
| --------------- | -------------------------------------------- | -------- |
| name            | <li>必填<li>長度大於1<li>長度小於100<li>唯一 | 服務名稱 |
| price           | <li>必填<li>數字最小是0                      | 價格     |
| durationMinutes | <li>必填<li>數字最小是0<li>小於1440          | 操作分鐘 |
| isAddon         | <li>必填<li>布林值                           | 附加服務 |
| isVisible       | <li>必填<li>布林值                           | 可見狀態 |
| note            | <li>選填<li>長度小於255                      | 備註     |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "9000000001",
    "name": "凝膠手部單色",
    "price": 1200,
    "durationMinutes": 60,
    "isAddon": false,
    "isVisible": true,
    "isActive": true,
    "note": "含基礎修型保養"
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "name": "name為必填項目"
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
  "message": "權限不足，僅限管理員新增服務"
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

1. 驗證 `name` 是否唯一。
2. 建立 `services` 資料。
3. 回傳新增結果。

---

## 注意事項

- 服務名稱不可重複。
- 支援「主服務」及「附加服務」類型。
- 設定為不可見或未啟用時，客戶不可從前台預約。
