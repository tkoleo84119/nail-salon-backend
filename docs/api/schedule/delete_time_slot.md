## User Story

1. 作為一位員工（`ADMIN` / `MANAGER` / `STYLIST`，不含 `SUPER_ADMIN`），我希望可以針對一筆 time_slot 進行刪除。
2. 作為一位管理員（`SUPER_ADMIN` / `ADMIN` / `MANAGER`），我希望可以針對單一美甲師的一筆 time_slot 進行刪除。

---

## Endpoint

**DELETE** `/api/schedules/{scheduleId}/time-slots/{timeSlotId}`

---

## 說明

- 可針對單一時段（time_slot）進行刪除。
- 員工僅能刪除自己班表的 time_slot（僅限自己有權限的 store）。
- 管理員可刪除任一美甲師的 time_slot（僅限自己有權限的 store）。
- 已被預約的 time_slot 禁止刪除。

---

## 權限

- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可刪除任何美甲師的 time_slot（僅限自己有權限的 store）。
- `STYLIST` 僅可刪除自己班表的 time_slot（僅限自己有權限的 store）。

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
| scheduleId | 班表ID |
| timeSlotId | 時段ID |

---

## Response

### 成功 204 No Content

```json
{
  "data": {
    "deleted": ["5000000011"]
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

#### 403 Forbidden - 權限不足/無法操作他人資料

```json
{
  "message": "權限不足，無法操作他人時段"
}
```

#### 404 Not Found - 時段不存在

```json
{
  "message": "時段不存在或已被刪除"
}
```

#### 409 Conflict - 已被預約之時段不可刪除

```json
{
  "message": "該時段已被預約，無法刪除"
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

- `time_slots`
- `schedules`
- `stylists`
- `stores`

---

## Service 邏輯

1. 檢查 `scheduleId` 是否存在。
2. 檢查 `timeSlotId` 是否存在。
3. 檢查 `timeSlotId` 是否屬於 `scheduleId`。
4. 判斷身分是否可操作指定 time_slot（員工僅可刪除自己的 time_slot，管理員可刪除任一美甲師 time_slot）。
5. 檢查是否有權限操作該 store。
6. 檢查 `timeSlotId` 是否已被預約。(被預約時，isAvailable = false，不可刪除)
7. 執行刪除。
8. 回傳已刪除 id。

---

## 注意事項

- 員工僅能針對自己的 time_slot 刪除；管理員可針對任一美甲師的 time_slot 刪除。
- 僅可操作自己有權限的 store。
- 時段一旦被預約（有 booking 記錄）禁止刪除。

