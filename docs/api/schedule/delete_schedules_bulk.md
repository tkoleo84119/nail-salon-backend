## User Story

1. 作為一位員工（`ADMIN` / `MANAGER` / `STYLIST`，不含 `SUPER_ADMIN`），我希望可以刪除我自己的多筆出勤班表（schedules）。
2. 作為一位管理員（`SUPER_ADMIN` / `ADMIN`），我希望可以刪除單一美甲師的多筆班表（schedules）。

---

## Endpoint

**DELETE** `/api/admin/schedules/bulk`

---

## 說明

- 一次只能針對同一位美甲師、同一家門市，刪除多筆班表（schedules）。
- 員工只能刪除自己班表。 （只能刪除自己有權限的 `store`）
- 管理員可刪除任一美甲師的班表。 （只能刪除自己有權限的 `store`）
- 刪除班表會一併刪除底下的 time_slots。

---

## 權限

- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可刪除任何美甲師的多筆班表。（只能刪除自己有權限的 `store`）
- `STYLIST` 僅可刪除自己班表。（只能刪除自己有權限的 `store`）

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Body（一次刪多筆班表）

```json
{
  "stylistId": "18000000001",
  "storeId": "1",
  "scheduleIds": ["4000000001", "4000000002"]
}
```

### 驗證規則

| 欄位        | 規則                         | 說明         |
| ----------- | ---------------------------- | ------------ |
| stylistId   | <li>必填                     | 美甲師id     |
| storeId     | <li>必填                     | 門市id       |
| scheduleIds | <li>必填<li>陣列<li>至少一筆 | 欲刪的班表id |

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

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "scheduleIds": "scheduleIds最小值為1"
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

#### 404 Not Found - 班表不存在

```json
{
  "message": "部分班表不存在或已被刪除"
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
2. 判斷身分是否可操作指定 stylistId (員工只能刪除自己的班表，管理員可刪除任一美甲師班表)。
3. 檢查 `storeId` 是否存在。
4. 判斷是否有權限操作指定 `storeId`。
5. 取得 `scheduleIds` 的班表資料（含底下所有 `time_slots`）。
6. 驗證 `scheduleIds` 是否屬於 `stylistId`/`storeId`。
7. 驗證 `scheduleIds` 的班表是否已被預約。
8. 執行刪除（含底下所有 `time_slots`）。
9. 回傳已刪除班表 id 陣列。

---

## 注意事項

- 員工僅能刪除自己班表；管理員可刪除任一美甲師班表。
- 刪除時應同時移除對應 `time_slots`。
