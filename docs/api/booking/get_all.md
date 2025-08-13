## User Story

作為顧客，我希望能夠取得我的預約資訊，並可依狀態分頁查詢，方便管理所有預約。

---

## Endpoint

**GET** `/api/bookings`

---

## 說明

- 提供顧客查詢自己所有的預約資訊。
- 支援分頁（limit、offset）。
- 支援排序（sort）。
- 不回傳 `status` 為 `NO_SHOW` 的預約。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Query Parameter

| 參數   | 型別   | 預設值 | 說明                                             |
| ------ | ------ | ------ | ------------------------------------------------ |
| limit  | int    | 20     | 單頁筆數                                         |
| offset | int    | 0      | 起始筆數                                         |
| sort   | string | -date  | 排序欄位 (可以逗號串接，有 `-` 表示 `DESC` 排序) |
| status | string |        | 預約狀態 (可多選，用逗號分隔)                    |

### 驗證規則

| 欄位   | 必填 | 其他規則                                          |
| ------ | ---- | ------------------------------------------------- |
| limit  | 否   | <li>最小值1<li>最大值100                          |
| offset | 否   | <li>最小值0<li>最大值1000000                      |
| sort   | 否   | <li>可以為 status, date (其餘會忽略)              |
| status | 否   | <li>值只能為`SCHEDULED`, `CANCELLED`, `COMPLETED` |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 10,
    "items": [
      {
        "id": "5000000001",
        "storeId": "8000000001",
        "storeName": "大安旗艦店",
        "stylistId": "2000000001",
        "stylistName": "Ava",
        "date": "2025-08-02",
        "timeSlotId": "3000000001",
        "startTime": "10:00",
        "endTime": "12:00",
        "status": "SCHEDULED"
      },
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

| 狀態碼 | 錯誤碼 | 常數名稱               | 說明                              |
| ------ | ------ | ---------------------- | --------------------------------- |
| 401    | E1002  | AuthInvalidCredentials | 無效的 accessToken，請重新登入    |
| 401    | E1003  | AuthTokenMissing       | accessToken 缺失，請重新登入      |
| 401    | E1004  | AuthTokenFormatError   | accessToken 格式錯誤，請重新登入  |
| 401    | E1006  | AuthContextMissing     | 未找到使用者認證資訊，請重新登入  |
| 401    | E1011  | AuthCustomerFailed     | 未找到有效的顧客資訊，請重新登入  |
| 400    | E2023  | ValFieldMinNumber      | {field} 最小值為 {param}          |
| 400    | E2026  | ValFieldMaxNumber      | {field} 最大值為 {param}          |
| 400    | E2027  | ValFieldOneof          | {field} 必須是 {param} 其中一個值 |
| 500    | E9001  | SysInternalError       | 系統發生錯誤，請稍後再試          |
| 500    | E9002  | SysDatabaseError       | 資料庫操作失敗                    |

---

## 資料表

- `bookings`
- `schedules`
- `time_slots`
- `stores`
- `stylists`

---

## Service 邏輯

1. 查詢 `customer_id` 的預約資訊。
2. 關聯查詢對應時段、美甲師、門市資訊 (排除 `status` 為 `NO_SHOW` 的預約)。
3. 回傳分頁結果。

---

## 注意事項

- 僅允許本人查詢。
- 狀態可多選（如 `status=SCHEDULED,COMPLETED`）。
- 順序依預約日期排序(DESC)。
