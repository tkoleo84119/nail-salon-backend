## User Story

作為員工，我希望可以查詢特定的排班資料，並一併取得該排班對應的美甲師資訊與所有時段（Time Slots），以便確認排班內容與剩餘時段。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/schedules/{scheduleId}`

---

## 說明

- 所有登入員工皆可查詢。
- 回傳指定排班資料與對應的美甲師與時段資訊。

---

## 權限

- 任一已登入員工皆可使用（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameters

| 參數       | 說明    |
| ---------- | ------- |
| storeId    | 門市 ID |
| scheduleId | 排班 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
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
}
```

### 失敗

#### 401 Unauthorized - 未登入

```json
{
  "message": "無效的 accessToken"
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
2. 驗證員工是否有權限存取該門市。
3. 驗證排班是否存在。
4. 查詢指定 `scheduleId` 是否存在，且 `store_id=storeId`。
5. JOIN `stylists` 表取得美甲師資料。
6. 查詢該筆排班底下所有 `time_slots`，依 `start_time` 排序。
7. 回傳合併資訊。

---

## 注意事項

- 一律回傳所有時段。
