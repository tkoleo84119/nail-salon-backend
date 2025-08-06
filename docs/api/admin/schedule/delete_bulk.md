## User Story

1. 作為一位美甲師，我希望可以刪除我自己的多筆出勤班表（schedules）。
2. 作為一位管理員，我希望可以刪除單一美甲師的多筆班表（schedules）。

---

## Endpoint

**DELETE** `/api/admin/store/:storeId/schedules/bulk`

---

## 說明

- 一次只能針對同一位美甲師、同一家門市，刪除多筆班表（schedules）。
- 美甲師只能刪除自己班表。 （只能刪除自己有權限的 `store`）
- 管理員可刪除任一美甲師的班表。（只能刪除自己有權限的 `store`）
- 刪除班表會一併刪除底下的 time_slots。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可為任何美甲師刪除班表 (只能刪除自己有權限的 `store`)。
- `STYLIST` 僅可為自己刪除班表 (只能刪除自己有權限的 `store`)。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "stylistId": "18000000001",
  "scheduleIds": ["4000000001", "4000000002"]
}
```

### 驗證規則

| 欄位        | 必填 | 其他規則                | 說明         |
| ----------- | ---- | ----------------------- | ------------ |
| stylistId   | 是   |                         | 美甲師id     |
| scheduleIds | 是   | <li>最小1筆<li>最大31筆 | 欲刪的班表id |

---

## Response

### 成功 204 No Content

```json
{
  "data": {
    "deleted": ["4000000001", "4000000002"]
  }
}
```

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                              |
| ------ | -------- | --------------------------------- |
| 401    | E1002    | 無效的 accessToken，請重新登入    |
| 401    | E1003    | accessToken 缺失，請重新登入      |
| 401    | E1004    | accessToken 格式錯誤，請重新登入  |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入  |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入  |
| 403    | E3001    | 權限不足，無法執行此操作          |
| 400    | E2002    | 路徑參數缺失，請檢查              |
| 400    | E2004    | 參數類型轉換失敗                  |
| 400    | E2020    | {field} 為必填項目                |
| 400    | E2022    | {field} 至少需要 {param} 個項目   |
| 400    | E2025    | {field} 最多只能有 {param} 個項目 |
| 400    | E3STO001 | 門市未啟用                        |
| 400    | E3SCH001 | 部分班表已被預約，無法刪除        |
| 400    | E3SCH003 | 部分班表不屬於指定的門市          |
| 400    | E3SCH004 | 部分班表不屬於指定的美甲師        |
| 404    | E3STO002 | 門市不存在或已被刪除              |
| 404    | E3STY001 | 美甲師資料不存在                  |
| 404    | E3SCH005 | 排班不存在或已被刪除              |
| 500    | E9001    | 系統發生錯誤，請稍後再試          |
| 500    | E9002    | 資料庫操作失敗                    |

#### 400 Bad Request - 參數驗證失敗

```json
{
  "errors": [
    {
      "code": "E2020",
      "message": "stylistId 為必填項目",
      "field": "stylistId"
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
2. 判斷身分是否可操作指定 stylistId (員工只能刪除自己的班表，管理員可刪除任一美甲師班表)。
3. 檢查 `storeId` 是否存在。
4. 判斷是否有權限操作指定 `storeId`。
5. 取得 `scheduleIds` 的班表資料（含底下所有 `time_slots`）。
6. 驗證 `scheduleIds` 是否屬於 `stylistId`/`storeId`。
7. 驗證 `scheduleIds` 的班表是否已被預約。
8. 執行刪除。
9. 回傳已刪除班表 id 陣列。

---

## 注意事項

- 員工僅能刪除自己班表；管理員可刪除任一美甲師班表。
- 刪除時應同時移除對應 `time_slots`。
