## User Story

作為顧客，我希望能夠取得某家店某位美甲師一段時間內的排班（Schedule），以便查詢可預約時段。

---

## Endpoint

**GET** `/api/stores/{storeId}/stylists/{stylistId}/schedules`

---

## 說明

- 支援已登入顧客查詢指定門市某位美甲師在一段期間內的排班。
- 每筆資料表示某位美甲師在哪一天尚有空檔。
- 若顧客為黑名單（`is_blacklisted=true`），回傳空陣列。
- 若美甲師無排班，則回傳空陣列。
- 依 `work_date` 升冪排序。

---

## 權限

- 僅顧客可查詢（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameter

| 參數      | 說明     |
| --------- | -------- |
| storeId   | 門市ID   |
| stylistId | 美甲師ID |

### Query Parameter

| 參數      | 型別   | 必填 | 說明                   |
| --------- | ------ | ---- | ---------------------- |
| startDate | string | 是   | 起始日期（YYYY-MM-DD） |
| endDate   | string | 是   | 結束日期（YYYY-MM-DD） |

- 期限不超過 60 天。

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 10,
    "items": [
      {"date": "2025-08-01", "available_slots": 3},
      {"date": "2025-08-02", "available_slots": 3},
      {"date": "2025-08-03", "available_slots": 2},
      {"date": "2025-08-05", "available_slots": 1},
      {"date": "2025-08-06", "available_slots": 3},
      {"date": "2025-08-07", "available_slots": 2},
      {"date": "2025-08-09", "available_slots": 3},
      {"date": "2025-08-10", "available_slots": 3},
    ]
  }
}
```

### 失敗

#### 400 Bad Request - 缺少參數

```json
{
  "message": "startDate 和 endDate 為必填欄位"
}
```

#### 404 Not Found - 門市不存在/未啟用

```json
{
  "message": "查無此門市或門市未啟用"
}
```

#### 404 Not Found - 美甲師不存在

```json
{
  "message": "查無此美甲師"
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
- `staff_users`
- `staff_user_store_access`
- `stylists`
- `schedules`
- `time_slots`
- `customers`

---

## Service 邏輯

1. 檢查顧客是否為黑名單（`is_blacklisted=true`），若是則回傳空陣列。
2. 驗證門市是否存在且為啟用狀態（`stores.is_active=true`）。
3. 若起始日期為過去，則將起始日期設為今天。
4. 查詢該門市該美甲師在指定日期範圍內的仍可預約的排班。
5. 回傳依 `work_date ASC` 排序之結果。

---

## 注意事項

- 回傳內容不包含時段（`time_slots`），僅為每日排班表。
- 欲查詢某排班的可預約時段，請使用另一 API。
- 查詢時間範圍不超過 60 天。