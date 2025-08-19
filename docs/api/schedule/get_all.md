## User Story

作為顧客，我希望能夠取得某家店某位美甲師一段時間內的排班（Schedule），以便查詢可預約時段。

---

## Endpoint

**GET** `/api/stores/{storeId}/stylists/{stylistId}/schedules`

---

## 說明

- 提供顧客查詢指定門市某位美甲師在一段期間內的排班。
- 每筆資料表示某位美甲師在哪一天尚有空檔，包含時段。
- 若美甲師無排班，則回傳空陣列。
- 若顧客為黑名單（`customers.is_blacklisted=true`），回傳空陣列。
- 依 `work_date` 升冪排序。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

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

### 驗證規則

| 欄位      | 必填 | 其他規則              |
| --------- | ---- | --------------------- |
| startDate | 是   | <li>格式為 YYYY-MM-DD |
| endDate   | 是   | <li>格式為 YYYY-MM-DD |

- 期限不超過 60 天。

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "schedules": [
      {"id": "1", "date": "2025-08-01"},
      {"id": "2", "date": "2025-08-02"},
      {"id": "3", "date": "2025-08-03"},
      {"id": "5", "date": "2025-08-05"},
      {"id": "6", "date": "2025-08-06"},
      {"id": "7", "date": "2025-08-07"},
      {"id": "9", "date": "2025-08-09"},
      {"id": "10", "date": "2025-08-10"},
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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                                                |
| ------ | ------ | ----------------------- | --------------------------------------------------- |
| 401    | E1002  | AuthInvalidCredentials  | 無效的 accessToken，請重新登入                      |
| 401    | E1003  | AuthTokenMissing        | accessToken 缺失，請重新登入                        |
| 401    | E1004  | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入                    |
| 401    | E1006  | AuthContextMissing      | 未找到使用者認證資訊，請重新登入                    |
| 401    | E1011  | AuthCustomerFailed      | 未找到有效的顧客資訊，請重新登入                    |
| 400    | E2002  | ValPathParamMissing     | 路徑參數缺失，請檢查                                |
| 400    | E2004  | ValTypeConversionFailed | 參數類型轉換失敗                                    |
| 400    | E2020  | ValFieldRequired        | {field} 為必填項目                                  |
| 400    | E2023  | ValFieldMinNumber       | {field} 最小值為 {param}                            |
| 400    | E2026  | ValFieldMaxNumber       | {field} 最大值為 {param}                            |
| 400    | E2033  | ValFieldDateFormat      | {field} 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD) |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試                            |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                                      |

---

## 資料表

- `stores`
- `staff_users`
- `staff_user_store_access`
- `stylists`
- `schedules`
- `time_slots`
- `customers`
- `bookings`

---

## Service 邏輯

1. 檢驗 endDate 是否在 startDate 之後。
1. 檢驗 endDate 與 startDate 之間的天數是否超過 31 天。
2. 檢查門市是否存在。
3. 檢查美甲師是否存在。
4. 檢查顧客是否為黑名單（`is_blacklisted=true`），若是則回傳空陣列。
5. 若起始日期為過去，則將起始日期設為今天。
6. 查詢該門市該美甲師在指定日期範圍內的仍可預約的排班。
7. 回傳依 `work_date ASC` 排序之結果。

---

## 注意事項

- 查詢時間範圍不超過 31 天。
- 回傳內容不包含時段（`time_slots`），僅為每日排班表。