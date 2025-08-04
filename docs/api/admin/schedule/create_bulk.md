## User Story

1. 作為一位美甲師，我希望可以安排我自己的出勤班表（schedules），每筆班表可包含多個 time_slot。
2. 作為一位管理員，我希望可以安排其他美甲師的班表（schedules），每筆班表可包含多個 time_slot。

---

## Endpoint

**POST** `/api/admin/store/:storeId/schedules/bulk`

---

## 說明

- 一次只能針對同一位美甲師、同一家門市，新增多日班表（schedules），每筆班表對應一個日期，可包含多個時段（time_slots）。
- 美甲師只能為自己建立班表 (只能建立自己有權限的 `store`)。
- 管理員可為任一美甲師建立班表 (只能建立自己有權限的 `store`)。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可為任何美甲師建立班表 (只能建立自己有權限的 `store`)。
- `STYLIST` 僅可為自己建立班表 (只能建立自己有權限的 `store`)。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameters

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Body 範例

```json
{
  "stylistId": "18000000001",
  "schedules": [
    {
      "workDate": "2024-07-21",
      "note": "早班",
      "timeSlots": [
        { "startTime": "09:00", "endTime": "12:00" },
        { "startTime": "13:00", "endTime": "18:00" }
      ]
    },
    {
      "workDate": "2024-07-22",
      "timeSlots": [
        { "startTime": "09:00", "endTime": "12:00" },
        { "startTime": "13:00", "endTime": "18:00" }
      ]
    }
  ]
}
```

### 驗證規則

| 欄位                          | 必填 | 其他規則                | 說明         |
| ----------------------------- | ---- | ----------------------- | ------------ |
| stylistId                     | 是   |                         | 美甲師id     |
| schedules                     | 是   | <li>最小1筆<li>最大31筆 | 多日班表     |
| schedules.workDate            | 是   | <li>YYYY-MM-DD 格式     | 班表日期     |
| schedules.note                | 否   | <li>最長100字元         | 備註         |
| schedules.timeSlots           | 是   | <li>最小1筆<li>最大20筆 | 當日多個時段 |
| schedules.timeSlots.startTime | 是   | <li>HH:mm 格式          | 起始時間     |
| schedules.timeSlots.endTime   | 是   | <li>HH:mm 格式          | 結束時間     |

---

## Response

### 成功 201 Created

```json
{
  "data": {
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
| 400    | E2022    | {field} 至少需要 {param} 個項目                     |
| 400    | E2024    | {field} 長度最多只能有 {param} 個字元               |
| 400    | E2025    | {field} 最多只能有 {param} 個項目                   |
| 400    | E2033    | {field} 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD) |
| 400    | E2034    | {field} 格式錯誤，請使用正確的時間格式 (HH:mm)      |
| 400    | E3SCH009 | 輸入的工作日期重複                                  |
| 400    | E3SCH010 | 不能創建過去的班表                                  |
| 400    | E3SCH011 | 時段時間區段重疊                                    |
| 400    | E3SCH012 | 結束時間必須在開始時間之後                          |
| 404    | E3STO002 | 門市不存在或已被刪除                                |
| 404    | E3STY001 | 美甲師資料不存在                                    |
| 409    | E3SCH013 | 美甲師班表已存在                                    |
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

- `schedules`
- `time_slots`
- `stylists`
- `stores`

---

## Service 邏輯

1. 檢查 `stylistId` 是否存在。
2. 判斷身分是否可操作指定 stylistId (員工只能建立自己的班表，管理員可建立任一美甲師班表)。
3. 檢查 `storeId` 是否存在。
4. 判斷是否有權限操作指定 `storeId`。
5. 驗證每筆 schedule 的 workDate、timeSlots
   - 驗證 workDate 格式是否正確。
   - 不可創建過去的班表。
   - 不可傳入相同的 workDate。
   - 驗證 timeSlots 的 startTime、endTime 格式是否正確。
   - 驗證 timeSlots 的 startTime 必須在 endTime 之前。
   - 驗證 timeSlots 的 startTime、endTime 不得重疊。
6. 檢查同一天同店同美甲師是否已有班表（不可重複排班）。
7. 新增 `schedules` 資料。
8. 批次建立對應的多筆 `time_slots`。
9. 回傳新增結果。

---

## 注意事項

- 員工僅能建立自己的班表；管理員可建立任一美甲師班表。
- 同一天、同店、同美甲師僅能有一筆 schedule。
- 每個 schedule 需至少一筆 time_slot，且時間區段不得重疊。
