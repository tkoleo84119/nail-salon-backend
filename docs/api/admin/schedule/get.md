## User Story

作為員工，我希望可以查詢特定的排班資料，並一併取得排班底下的時段（Time Slots），以便確認排班內容。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/schedules/{scheduleId}`

---

## 說明

- 回傳指定排班資料與對應的時段資訊。
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

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                             |
| ------ | -------- | -------------------------------- |
| 401    | E1002    | 無效的 accessToken，請重新登入   |
| 401    | E1003    | accessToken 缺失，請重新登入     |
| 401    | E1004    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入 |
| 403    | E3001    | 權限不足，無法執行此操作         |
| 400    | E2002    | 路徑參數缺失，請檢查             |
| 400    | E2004    | 參數類型轉換失敗                 |
| 400    | E2020    | {field} 為必填項目               |
| 404    | E3STO002 | 門市不存在或已被刪除             |
| 404    | E3SCH005 | 排班不存在或已被刪除             |
| 500    | E9001    | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | 資料庫操作失敗                   |


#### 400 Bad Request - 參數驗證失敗

```json
{
  "errors": [
    {
      "code": "E2002",
      "message": "路徑參數缺失，請檢查"
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

1. 驗證門市是否存在。
2. 驗證員工是否有權限存取該門市。
3. 查詢指定 `scheduleId` 同時 JOIN `time_slots` 表取得時段資料。
4. 回傳合併資訊。

---

## 注意事項

- 一律回傳所有時段。
