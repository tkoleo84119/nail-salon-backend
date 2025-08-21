## User Story

作為員工，我希望可以取消某筆顧客預約，並可填寫取消原因，以利記錄與管理預約異動。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}/bookings/{bookingId}/cancel`

---

## 說明

- 提供後台管理員取消預約功能。
- 可記錄取消原因。
- 狀態可以變更為 `CANCELLED` 或 `NO_SHOW`。
- 預約取消後，會釋放對應時段（`time_slots.is_available=true`）。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Authorization: Bearer <access_token>
- Content-Type: application/json

### Path Parameters

| 參數      | 說明    |
| --------- | ------- |
| storeId   | 門市 ID |
| bookingId | 預約 ID |

### Body 選填

```json
{
  "status": "CANCELLED",
  "cancelReason": "顧客臨時無法前來"
}
```

#### 驗證規則
| 欄位         | 必填 | 其他規則                       | 說明     |
| ------------ | ---- | ------------------------------ | -------- |
| status       | 是   | <li>值只能是CANCELLED或NO_SHOW | 取消狀態 |
| cancelReason | 否   | <li>最長255字                  | 取消原因 |

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

| 狀態碼 | 錯誤碼   | 常數名稱                   | 說明                                  |
| ------ | -------- | -------------------------- | ------------------------------------- |
| 401    | E1002    | AuthTokenInvalid           | 無效的 accessToken，請重新登入        |
| 401    | E1003    | AuthTokenMissing           | accessToken 缺失，請重新登入          |
| 401    | E1004    | AuthTokenFormatError       | accessToken 格式錯誤，請重新登入      |
| 401    | E1006    | AuthContextMissing         | 未找到使用者認證資訊，請重新登入      |
| 401    | E1011    | AuthCustomerFailed         | 未找到有效的顧客資訊，請重新登入      |
| 400    | E2001    | ValJSONFormatError         | JSON 格式錯誤，請檢查                 |
| 400    | E2024    | ValFieldStringMaxLength    | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2030    | ValFieldOneOf              | {field} 必須是 {param} 其中一個值     |
| 400    | E3SER001 | ServiceNotActive           | 服務未啟用                            |
| 400    | E3SER002 | ServiceNotMainService      | 服務不是主服務                        |
| 400    | E3SER003 | ServiceNotAddon            | 服務不是附屬服務                      |
| 400    | E3TMS006 | TimeSlotNotEnoughTime      | 時段時間不足                          |
| 404    | E3TMS005 | TimeSlotNotFound           | 時段不存在或已被刪除                  |
| 404    | E3SER004 | ServiceNotFound            | 服務不存在或已被刪除                  |
| 404    | E3STY001 | StylistNotFound            | 美甲師資料不存在                      |
| 409    | E3BK006  | BookingTimeSlotUnavailable | 該時段已被預約，請重新選擇            |
| 500    | E9001    | SysInternalError           | 系統發生錯誤，請稍後再試              |
| 500    | E9002    | SysDatabaseError           | 資料庫操作失敗                        |

---

## 資料表

- `bookings`
- `time_slots`

---

## Service 邏輯
1. 驗證預約是否存在
2. 檢查預約狀態是否為 SCHEDULED
3. 更新 `status` 並寫入 `cancel_reason`
4. 將該預約所屬 `time_slots.is_available = true`
5. 回傳更新後狀態

---

## 注意事項

- 僅允許取消尚未完成的預約。
