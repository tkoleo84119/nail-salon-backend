## User Story

作為員工，我希望可以查詢某門市下的所有排班資料，並一併取得每筆排班底下的時段（Time Slots），並支援查詢條件與分頁，以利安排預約與時段管理。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/schedules`

---

## 說明

- 所有登入員工皆可查詢。
- 回傳每筆排班資料（schedule），並包含該筆排班的時段（time slots）。
- 支援查詢條件：美甲師、日期區間、是否可預約時段（is_available）、分頁。

---

## 權限

- 任一已登入員工皆可使用（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameters

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Query Parameters

| 參數        | 型別   | 必填 | 預設值 | 說明                              |
| ----------- | ------ | ---- | ------ | --------------------------------- |
| stylistId   | string | 否   |        | 篩選指定美甲師的排班              |
| startDate   | string | 是   |        | 起始排班日期（YYYY-MM-DD）        |
| endDate     | string | 是   |        | 結束排班日期（YYYY-MM-DD）        |
| isAvailable | bool   | 否   |        | 是否還可預約（is_available=true） |

- startDate 與 endDate 間隔最多 60 天

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "items": [
      {
        "id": "5000000001",
        "workDate": "2025-08-01",
        "stylist": {
          "id": "7000000001",
          "name": "Ariel"
        },
        "note": "上班全天",
        "timeSlots": [
          {
            "id": "9000000001",
            "startTime": "10:00",
            "endTime": "11:00",
            "isAvailable": true
          },
          {
            "id": "9000000002",
            "startTime": "11:00",
            "endTime": "12:00",
            "isAvailable": false
          }
        ]
      }
    ]
  }
}
```

### 失敗

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 無權限

```json
{
  "message": "權限不足，無法執行此操作"
}
```

#### 404 Not Found - 門市不存在

```json
{
  "message": "門市不存在或已被刪除"
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
- `schedules`
- `stylists`
- `time_slots`

---

## Service 邏輯

1. 驗證門市是否存在。
2. 驗證員工是否擁有該門市存取權限。
3. 查詢該門市之 `schedules`，依條件過濾：
   - `stylist_id`, `work_date BETWEEN startDate AND endDate`
4. JOIN `stylists` 表取得名稱、staff_user_id。
5. JOIN `time_slots` 表取得該筆排班之所有時段，並可依 `isAvailable` 過濾。
6. 回傳分頁資料與總筆數。

---

## 注意事項

- 排班依 `work_date` 排序。
- 時段依 `start_time` 排序。
