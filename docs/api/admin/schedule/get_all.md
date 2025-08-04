## User Story

作為員工，我希望可以查詢某門市下的所有排班資料，並一併取得每筆排班底下的時段（Time Slots），以利安排預約與時段管理。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/schedules`

---

## 說明

- 回傳每筆排班資料（schedule），並包含該筆排班的時段（time slots）。
- 支援查詢條件：美甲師、日期區間。
- 排班依 `work_date` 排序。
- 時段依 `start_time` 排序。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameters

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Query Parameters

| 參數        | 型別   | 必填 | 預設值 | 說明                                         |
| ----------- | ------ | ---- | ------ | -------------------------------------------- |
| stylistId   | string | 否   |        | 篩選指定美甲師的排班(可傳入多個，以逗號分隔) |
| startDate   | string | 是   |        | 起始排班日期（YYYY-MM-DD）                   |
| endDate     | string | 是   |        | 結束排班日期（YYYY-MM-DD）                   |
| isAvailable | bool   | 否   |        | 是否可預約                                   |

### 驗證規則

| 欄位        | 必填 | 其他規則                |
| ----------- | ---- | ----------------------- |
| stylistId   | 否   |                         |
| startDate   | 是   | <li>格式要是 YYYY-MM-DD |
| endDate     | 是   | <li>格式要是 YYYY-MM-DD |
| isAvailable | 否   | <li>是否是布林值        |

- startDate 與 endDate 間隔最多 31 天

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "stylistList": [
      {
        "id": "7000000001",
        "name": "Ariel",
        "schedules": [
          {
            "id": "5000000001",
            "workDate": "2025-08-01",
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
    ]
  }
}
```

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                                                |
| ------ | -------- | --------------------------------------------------- |
| 401    | E1002    | 無效的 accessToken，請重新登入                      |
| 401    | E1003    | accessToken 缺失，請重新登入                        |
| 401    | E1004    | accessToken 格式錯誤，請重新登入                    |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入                    |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入                    |
| 403    | E3001    | 權限不足，無法執行此操作                            |
| 400    | E2002    | 路徑參數缺失，請檢查                                |
| 400    | E2004    | 參數類型轉換失敗                                    |
| 400    | E2020    | {field} 為必填項目                                  |
| 400    | E2029    | {field} 必須是布林值                                |
| 400    | E2033    | {field} 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD) |
| 400    | E3SCH007 | 結束日期必須在開始日期之後                          |
| 400    | E3SCH008 | 日期範圍不能超過 31 天                              |
| 404    | E3STO002 | 門市不存在或已被刪除                                |
| 500    | E9001    | 系統發生錯誤，請稍後再試                            |
| 500    | E9002    | 資料庫操作失敗                                      |

#### 400 Bad Request - 參數驗證失敗

```json
{
  "errors": [
    {
      "code": "E2020",
      "message": "startDate 為必填項目",
      "field": "startDate"
    }
  ]
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "errors": [
    {
      "code": "E1002",
      "message": "無效的 accessToken，請重新登入"
    }
  ]
}
```

#### 403 Forbidden - 無權限

```json
{
  "errors": [
    {
      "code": "E3001",
      "message": "權限不足，無法執行此操作"
    }
  ]
}
```

#### 404 Not Found - 門市不存在

```json
{
  "errors": [
    {
      "code": "E3STO002",
      "message": "門市不存在或已被刪除"
    }
  ]
}
```

#### 500 Internal Server Error - 系統錯誤

```json
{
  "errors": [
    {
      "code": "E9001",
      "message": "系統發生錯誤，請稍後再試"
    }
  ]
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

1. 驗證日期邏輯
   - 結束日期必須在開始日期之後
   - 日期範圍不能超過 31 天
2. 驗證門市是否存在。
3. 驗證員工是否擁有該門市存取權限。
4. 查詢該門市之 `schedules`，依條件過濾：
   - `stylist_id`, `work_date BETWEEN startDate AND endDate`
   - 加入 `isAvailable` 過濾條件
5. JOIN `stylists` 和 `time_slots` 表取得該筆排班之所有時段。
6. 整理排班資料，依美甲師分組。
7. 回傳排班資料。

---

## 注意事項

- 排班依 `work_date` 排序，時段依 `start_time` 排序。
