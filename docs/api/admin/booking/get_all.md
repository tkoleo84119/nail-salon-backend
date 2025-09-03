## User Story

作為員工，我希望可以查詢某門市下所有預約紀錄（Booking），並支援條件查詢與分頁，方便查看排程與歷史記錄。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/bookings`

---

## 說明

- 支援基本查詢條件。
- 支援分頁（limit、offset）。
- 支援排序（sort）。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Query Parameters

| 參數      | 型別   | 必填 | 預設值      | 說明                                                 |
| --------- | ------ | ---- | ----------- | ---------------------------------------------------- |
| stylistId | string | 否   |             | 篩選指定美甲師的預約                                 |
| startDate | string | 否   |             | 起始日期（YYYY-MM-DD）                               |
| endDate   | string | 否   |             | 結束日期（YYYY-MM-DD）                               |
| status    | string | 否   |             | 預約狀態（SCHEDULED, CANCELLED, COMPLETED, NO_SHOW） |
| limit     | int    | 否   | 20          | 單頁筆數                                             |
| offset    | int    | 否   | 0           | 起始筆數                                             |
| sort      | string | 否   | -created_at | 排序欄位 (可以逗號串接，有 `-` 表示 `DESC` 排序)     |

### 驗證規則
| 欄位      | 必填 | 其他規則                                                   |
| --------- | ---- | ---------------------------------------------------------- |
| stylistId | 否   |                                                            |
| startDate | 否   |                                                            |
| endDate   | 否   |                                                            |
| status    | 否   | <li>只能為 SCHEDULED, CANCELLED, COMPLETED, NO_SHOW        |
| limit     | 否   | <li>最小值1<li>最大值100                                   |
| offset    | 否   | <li>最小值0<li>最大值1000000                               |
| sort      | 否   | <li>可以為 createdAt, updatedAt, status, date (其餘會忽略) |

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
        "actualDuration": 120, // 如果非 "COMPLETED" 則不會有該欄位
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

### 錯誤處理

全部 API 皆回傳如下結構，請參考錯誤總覽。

```json
{
  "errors": [
    {
      "code": "EXXXX",
      "message": "錯誤訊息",
      "field": "錯誤欄位名稱"
    }
  ]
}
```

- 欄位說明：
  - errors: 錯誤陣列（支援多筆同時回報）
  - code: 錯誤代碼，唯一對應每種錯誤
  - message: 中文錯誤訊息（可參照錯誤總覽）
  - field: 參數欄位名稱（僅部分驗證錯誤有）

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                              |
| ------ | -------- | ----------------------- | --------------------------------- |
| 401    | E1002    | AuthTokenInvalid        | 無效的 accessToken，請重新登入    |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入      |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入  |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入  |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入  |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作          |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查              |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                  |
| 400    | E2023    | ValFieldMinNumber       | {field} 最小值為 {param}          |
| 400    | E2026    | ValFieldMaxNumber       | {field} 最大值為 {param}          |
| 400    | E2030    | ValFieldOneof           | {field} 必須是 {param} 其中一個值 |
| 404    | E3STO002 | StoreNotFound           | 門市不存在或已被刪除              |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試          |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                    |

---

## 資料表

- `bookings`
- `customers`
- `stylists`
- `time_slots`
- `schedules`
- `booking_details`
- `services`

---

## Service 邏輯

1. 驗證 `storeId` 是否存在。
2. 驗證員工是否有權限查詢該門市。
3. 查詢 `bookings` 資料
4. 加入 `limit` / `offset` 分頁，和 `sort` 排序
5. 查詢 `services` 資料
6. 整理資料，並依 `services` 的 `is_addon` 為 `true` 則為子服務，反之為主服務。
