## User Story

1. 作為一位員工（`ADMIN` / `MANAGER` / `STYLIST`，不含 `SUPER_ADMIN`），我希望可以針對一筆 time_slot 進行更新（時間與是否可被預約）。
2. 作為一位管理員（`SUPER_ADMIN` / `ADMIN` / `MANAGER`），我希望可以針對單一美甲師的一筆 time_slot 進行更新（時間與是否可被預約）。

---

## Endpoint

**PATCH** `/api/schedules/{scheduleId}/time-slots/{timeSlotId}`

---

## 說明

- 可針對單一時段（time_slot）進行更新，包括起訖時間、是否可預約。
- 員工僅能編輯自己班表的 time_slot（僅限自己有權限的 store）。
- 管理員可編輯任一美甲師的 time_slot（僅限自己有權限的 store）。
- 被預約時段不可更新。
- 要更新時段必須同時傳入 startTime/endTime 兩個欄位，不可單獨傳入。
- 更新時需檢查時間區段是否與同一 schedule 下其他 time_slots 重疊。

---

## 權限

- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可更新任何美甲師的 time_slot（僅限自己有權限的 store）。
- `STYLIST` 僅可更新自己班表的 time_slot（僅限自己有權限的 store）。

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

### Body

```json
{
  "startTime": "14:00",
  "endTime": "16:00",
  "isAvailable": true
}
```

### 驗證規則

| 欄位        | 規則                   | 說明         |
| ----------- | ---------------------- | ------------ |
| startTime   | <li>選填<li>HH:mm 格式 | 時段起始時間 |
| endTime     | <li>選填<li>HH:mm 格式 | 時段結束時間 |
| isAvailable | <li>選填               | 是否可被預約 |

- 至少需有一項欄位出現

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "5000000011",
    "scheduleId": "4000000001",
    "startTime": "14:00",
    "endTime": "16:00",
    "isAvailable": true
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤 / 時段重疊

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "startTime": "startTime格式錯誤"
  }
}
```

```json
{
  "message": "時段重疊，不可設定重複時間"
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

#### 409 Conflict - 時段重疊

```json
{
  "message": "時段重疊，不可設定重複時間"
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

1. 驗證至少一個欄位有更新。
2. 若 startTime/endTime 有傳入，另外一個欄位必須有傳入。
3. 檢查 `timeSlotId` 是否存在。
4. 判斷 time_slot 是否屬於指定 schedule。
5. 檢查 `timeSlotId` 是否已被預約。(被預約時，不可變更任何欄位)
6. 檢查 `scheduleId` 是否存在。
7. 取得 stylist 資訊。
8. 判斷身分是否可操作指定 time_slot（員工僅可編輯自己的 time_slot，管理員可編輯任一美甲師 time_slot）。
9.  檢查是否有權限操作該 store。
10. 若有更新時間，檢查是否時間相關邏輯
    1.  startTime / endTime 格式是否正確。
    2.  startTime 必須在 endTime 之前。
    3.  startTime / endTime 是否與 schedule 下其他 time_slots 重疊。
11. 更新 time_slot。
12. 回傳更新結果。

---

## 注意事項

- 員工僅能針對自己的 time_slot 編輯；管理員可針對任一美甲師的 time_slot 編輯。
- 僅可操作自己有權限的 store。
- 若只更新 isAvailable，則時間不檢查重疊；若有更動時間則必須比對同 schedule 下其他 time_slots 是否重疊。
- 被預約時段不可更新。
