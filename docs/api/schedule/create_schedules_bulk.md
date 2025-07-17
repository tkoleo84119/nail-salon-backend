## User Story

1. 作為一位員工（`ADMIN` / `MANAGER` / `STYLIST`，不含 `SUPER_ADMIN`），我希望可以安排我自己的出勤班表（schedules），每筆班表可包含多個 time_slot。
2. 作為一位管理員（`SUPER_ADMIN` / `ADMIN` / `MANAGER`），我希望可以安排其他美甲師的班表（schedules），每筆班表可包含多個 time_slot。

---

## Endpoint

**POST** `/api/schedules/bulk`

---

## 說明

- 一次只能針對同一位美甲師、同一家門市，新增多日班表（schedules），每筆班表對應一個日期，可包含多個時段（time_slots）。
- 員工只能為自己建立班表 (只能建立自己有權限的 `store`)。
- 管理員可為任一美甲師建立班表 (只能建立自己有權限的 `store`)。

---

## 權限

- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可為任何美甲師建立班表 (只能建立自己有權限的 `store`)。
- `STYLIST` 僅可為自己建立班表 (只能建立自己有權限的 `store`)。

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Body（一次可新增多日班表）

```json
{
  "stylistId": "18000000001",
  "storeId": "1",
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

| 欄位                          | 規則                        | 說明         |
| ----------------------------- | --------------------------- | ------------ |
| stylistId                     | <li>必填                    | 美甲師id     |
| storeId                       | <li>必填                    | 門市id       |
| schedules                     | <li>必填<li>陣列，不可為空  | 多日班表     |
| schedules.workDate            | <li>必填<li>YYYY-MM-DD 格式 | 班表日期     |
| schedules.note                | <li>選填<li>長度≤100        | 備註         |
| schedules.timeSlots           | <li>必填<li>陣列            | 當日多個時段 |
| schedules.timeSlots.startTime | <li>必填<li>HH:mm 格式      | 起始時間     |
| schedules.timeSlots.endTime   | <li>必填<li>HH:mm 格式      | 結束時間     |

---

## Response

### 成功 201 Created

```json
{
  "data": [
    {
      "scheduleId": "4000000001",
      "stylistId": "18000000001",
      "storeId": "1",
      "workDate": "2024-07-21",
      "note": "早班",
      "timeSlots": [
        { "id": "5000000001", "startTime": "09:00", "endTime": "12:00" },
        { "id": "5000000002", "startTime": "13:00", "endTime": "18:00" }
      ]
    }
  ]
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "schedules[0].workDate": "格式錯誤",
    "schedules[1].timeSlots": "不得為空"
  }
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 權限不足/無法操作他人資料

```json
{
  "message": "權限不足，無法操作他人班表"
}
```

#### 409 Conflict - 已存在班表

```json
{
  "message": "美甲師班表已存在"
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
5. 驗證每筆 schedule 的日期、timeSlots（不可重複）。
6. 檢查同一天同店同美甲師是否已有班表（不可重複排班）。
7. 新增 `schedules` 資料。
8. 批次建立對應的多筆 `time_slots`。
9. 回傳所有新增結果。

---

## 注意事項

- 員工僅能建立自己的班表；管理員可建立任一美甲師班表。
- 同一天、同店、同美甲師僅能有一筆 schedule。
- 每個 schedule 需至少一筆 time_slot，且時間區段不得重疊。
