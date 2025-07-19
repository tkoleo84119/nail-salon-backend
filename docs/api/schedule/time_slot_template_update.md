## User Story

作為一位管理員（`SUPER_ADMIN` / `ADMIN` / `MANAGER`），我希望能更新時段範本（template）的名稱與備註，方便後續管理。

---

## Endpoint

**PATCH** `/api/time-slot-templates/{templateId}`

---

## 說明

- 僅限管理員可更新指定範本（template）的 name 與 note。
- 僅支援名稱與備註修改，範本下的時段（timeSlots）不在本 API 內調整。

---

## 權限

- 僅 `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可操作。

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Path Parameter

| 參數       | 說明   |
| ---------- | ------ |
| templateId | 範本ID |

### Body

```json
{
  "name": "新標準早班",
  "note": "夏季適用"
}
```

- 兩者皆為選填，但至少要有一項。

### 驗證規則

| 欄位 | 規則                                | 說明     |
| ---- | ----------------------------------- | -------- |
| name | <li>選填<li>長度大於1<li>長度小於50 | 範本名稱 |
| note | <li>選填<li>長度小於100             | 備註     |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "6000000011",
    "name": "新標準早班",
    "note": "夏季適用"
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "name": "name長度不可超過50字"
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

#### 404 Not Found - 範本不存在

```json
{
  "message": "範本不存在或已被刪除"
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

- `time_slot_templates`

---

## Service 邏輯

1. 驗證至少一個欄位有更新。
2. 驗證 `templateId` 是否存在。
3. 更新 `time_slot_templates` 資料。
4. 回傳更新結果。

---

## 注意事項

- 僅管理員可操作。
- 僅允許 name 與 note 欄位修改。

