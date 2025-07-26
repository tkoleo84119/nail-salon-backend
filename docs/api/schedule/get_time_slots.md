## User Story

作為顧客，我希望能夠取得某筆排班（Schedule）底下仍可以預約的時段（Time Slot），以便選擇預約時間。

---

## Endpoint

**GET** `/api/schedules/{scheduleId}/time-slots`

---

## 說明

- 支援已登入顧客查詢指定排班（schedule）下可預約的時段。
- 僅回傳 `is_available=true` 的時段。
- 每筆資料包含：起始時間、結束時間與該時段的總分鐘數（供前端計算服務是否可容納）。
- 若顧客為黑名單（`is_blacklisted=true`），回傳空陣列。
- 若排班不存在，回傳空陣列。
- 按照 `startTime` Asc 排序。

---

## 權限

- 僅顧客可查詢（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameter

| 參數       | 說明    |
| ---------- | ------- |
| scheduleId | 排班 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": [
    {
      "id": "9000000001",
      "startTime": "10:00",
      "endTime": "11:00",
      "durationMinutes": 60
    },
    {
      "id": "9000000002",
      "startTime": "11:30",
      "endTime": "12:30",
      "durationMinutes": 60
    }
  ]
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
- `schedules`
- `time_slots`

---

## Service 邏輯

1. 檢查顧客是否為黑名單（`is_blacklisted=true`），若是則回傳空陣列。
2. 驗證指定排班是否存在，若不存在則回傳空陣列。
3. 查詢該排班底下 `is_available=true` 的時段，並計算 `durationMinutes = endTime - startTime`（單位：分鐘）。
4. 回傳每筆時段。

---

## 注意事項

- `durationMinutes` 為方便前端驗證服務長度是否合適。
