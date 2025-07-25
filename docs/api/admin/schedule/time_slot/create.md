## User Story

1. 作為一位員工（`ADMIN` / `MANAGER` / `STYLIST`，不含 `SUPER_ADMIN`），我希望可以針對自己的班表（schedule）新增一個時段（time_slot）。
2. 作為一位管理員（`SUPER_ADMIN` / `ADMIN` / `MANAGER`），我希望可以針對單一美甲師的班表（schedule）新增一個時段（time_slot）。

---

## Endpoint

**POST** `/api/admin/schedules/{scheduleId}/time-slots`

---

## 說明

- 可針對單一班表（schedule）新增一個時段（time_slot）。
- 員工只能針對自己的班表新增時段（僅限自己有權限的 store）。
- 管理員可針對任一美甲師的班表新增時段（僅限自己有權限的 store）。
- 新增時段時，需檢查時段是否與既有時段重疊。

---

## 權限

- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可為任何美甲師的班表新增時段（僅限自己有權限的 store）。
- `STYLIST` 僅可為自己的班表新增時段（僅限自己有權限的 store）。

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Path Parameter

| 參數       | 說明   |
| ---------- | ------ |
| scheduleId | 班表ID |

### Body

```json
{
  "startTime": "14:00",
  "endTime": "16:00"
}
```

### 驗證規則

| 欄位      | 規則                   | 說明         |
| --------- | ---------------------- | ------------ |
| startTime | <li>必填<li>HH:mm 格式 | 時段起始時間 |
| endTime   | <li>必填<li>HH:mm 格式 | 時段結束時間 |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "5000000011",
    "scheduleId": "4000000001",
    "startTime": "14:00",
    "endTime": "16:00",
    "isAvailable": true
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤 / 時段重疊

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "startTime": "startTime為必填項目"
  }
}
```

```json
{
  "message": "時段重疊，不可新增重複時間"
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

#### 404 Not Found - 班表不存在

```json
{
  "message": "班表不存在或已被刪除"
}
```

#### 409 Conflict - 時段重疊

```json
{
  "message": "時間區段重疊"
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

1. 檢查 startTime/endTime 格式。
2. 確認 startTime 必須在 endTime 之前。
3. 檢查 `scheduleId` 是否存在。
4. 檢查 schedule 所屬的 stylist/store 是否存在。
5. 判斷身分是否可操作指定 schedule（員工只能新增自己的班表，管理員可新增任一美甲師班表）。
6. 檢查是否有權限操作該 store。
7. 檢查時間區間是否與該 schedule 下既有 time_slots 重疊。
8. 建立新的 time_slot，並設定 isAvailable 為 true。
9. 回傳新增結果。

---

## 注意事項

- 員工僅能針對自己的班表新增時段；管理員可針對任一美甲師的班表新增。
- 不可新增重疊時段（需比對現有 time_slots）。
- 僅可操作自己有權限的 store。

