## User Story

作為員工，我希望可以查詢某門市下所有預約紀錄（Booking），並支援條件查詢與分頁，方便查看排程與歷史記錄。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/bookings`

---

## 說明

- 所有登入員工皆可查詢。
- 支援依美甲師、起訖日篩選，並支援分頁。
- 回傳依 `workDate` 升冪排序。

---

## 權限

- 任一已登入員工皆可使用（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Query Parameters

| 參數      | 型別   | 必填 | 預設值 | 說明                   |
| --------- | ------ | ---- | ------ | ---------------------- |
| stylistId | string | 否   |        | 篩選指定美甲師的預約   |
| startDate | string | 否   |        | 起始日期（YYYY-MM-DD） |
| endDate   | string | 否   |        | 結束日期（YYYY-MM-DD） |
| limit     | int    | 否   | 20     | 單頁筆數               |
| offset    | int    | 否   | 0      | 起始筆數               |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 2,
    "items": [
      {
        "id": "3000000001",
        "customer": {
          "id": "2000000001",
          "name": "小美"
        },
        "stylist": {
          "id": "7000000001",
          "name": "Ariel"
        },
        "timeSlot": {
          "id": "9000000001",
          "workDate": "2025-08-01",
          "startTime": "10:00",
          "endTime": "11:00"
        },
        "mainService": {
          "id": "9000000010",
          "name": "法式美甲"
        },
        "subServices": [
          {
            "id": "9000000012",
            "name": "跳色"
          }
        ],
        "status": "SCHEDULED"
      }
    ]
  }
}
```

### 失敗

#### 401 Unauthorized

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

- `bookings`
- `customers`
- `stylists`
- `time_slots`
- `schedules`

---

## Service 邏輯

1. 驗證 `storeId` 是否存在。
2. 驗證員工是否有權限查詢該門市。
3. JOIN `bookings` → `time_slots` → `schedules` → `work_date`, `stylist_id`。
4. 加入查詢條件：
   - `stylistId` 篩選對應美甲師
   - `work_date BETWEEN startDate AND endDate`
5. 加入 `limit` / `offset` 分頁，並依 `work_date` 升冪排序。
6. JOIN `customers`, `stylists` 補足顯示欄位。

---

