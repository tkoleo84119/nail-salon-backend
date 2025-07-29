## User Story

作為員工，我希望可以取消某筆顧客預約，並可填寫取消原因，以利記錄與管理預約異動。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}/bookings/{bookingId}/cancel`

---

## 說明

- 所有登入員工皆可操作。
- 用於取消某筆預約，並可記錄取消原因。
- 狀態可以變更為 `CANCELLED` 或 `NO_SHOW`。
- 預約取消後，會釋放對應時段（`time_slots.is_available=true`）。

---

## 權限

- 任一已登入員工皆可使用（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>
Content-Type: application/json

### Path Parameters

| 參數      | 說明    |
| --------- | ------- |
| storeId   | 門市 ID |
| bookingId | 預約 ID |

### Body（可選）

```json
{
  "status": "CANCELLED",
  "cancelReason": "顧客臨時無法前來"
}
```

#### 驗證規則
| 欄位         | 規則                             | 說明     |
| ------------ | -------------------------------- | -------- |
| status       | <li>必填<li>CANCELLED 或 NO_SHOW | 取消狀態 |
| cancelReason | <li>選填<li>最長100字            | 取消原因 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "3000000001"
  }
}
```

### 失敗

#### 404 Not Found - 預約不存在

```json
{
  "message": "查無此門市或預約"
}
```

#### 409 Conflict - 已完成或已取消

```json
{
  "message": "該預約已完成或已取消，無法重複操作"
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

- `bookings`
- `time_slots`

---

## Service 邏輯

1. 驗證門市是否存在
2. 驗證員工是否有權限操作該門市
3. 驗證預約是否存在
4. 檢查預約狀態是否為 SCHEDULED
5. 更新 `status` 並寫入 `cancel_reason`
6. 將該預約所屬 `time_slots.is_available = true`
7. 回傳更新後狀態

---

## 注意事項

- 僅允許取消尚未完成的預約。

