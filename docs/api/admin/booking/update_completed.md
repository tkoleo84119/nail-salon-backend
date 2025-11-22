## User Story

作為員工，我希望可以更新某筆顧客已完成預約（Booking）的資料，如：完成時間。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}/bookings/{bookingId}/completed`

---

## 說明

- 提供後台管理員更新預約完成時間。
- 可修改指定預約的內容：完成時間。

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

### Body 範例

```json
{
  "actualDuration": 120,
  "pinterestImageUrls": [
    "https://pin.it/xxxxx",
    "https://pin.it/xxxxx"
  ]
}
```

### 驗證規則

| 欄位               | 必填 | 其他規則                  | 說明               |
| ------------------ | ---- | ------------------------- | ------------------ |
| actualDuration     | 否   | <li>最小值0<li>最大值1440 | 完成時間(分)       |
| pinterestImageUrls | 否   | <li>最多5張               | Pinterest 圖片 URL |


- 至少提供一個欄位進行更新

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

| 狀態碼 | 錯誤碼  | 常數名稱                        | 說明                             |
| ------ | ------- | ------------------------------- | -------------------------------- |
| 401    | E1002   | AuthTokenInvalid                | 無效的 accessToken，請重新登入   |
| 401    | E1003   | AuthTokenMissing                | accessToken 缺失，請重新登入     |
| 401    | E1004   | AuthTokenFormatError            | accessToken 格式錯誤，請重新登入 |
| 401    | E1006   | AuthContextMissing              | 未找到使用者認證資訊，請重新登入 |
| 401    | E1011   | AuthCustomerFailed              | 未找到有效的顧客資訊，請重新登入 |
| 403    | E1010   | AuthPermissionDenied            | 權限不足，無法執行此操作         |
| 400    | E2001   | ValJSONFormatError              | JSON 格式錯誤，請檢查            |
| 400    | E2002   | ValPathParamMissing             | 路徑參數缺失，請檢查             |
| 400    | E2003   | ValAllFieldsEmpty               | 至少需要提供一個欄位進行更新     |
| 400    | E2004   | ValTypeConversionFailed         | 參數類型轉換失敗                 |
| 400    | E2023   | ValFieldMinNumber               | {field} 最小值為 {param}         |
| 400    | E2026   | ValFieldMaxNumber               | {field} 最大值為 {param}         |
| 400    | E3BK002 | BookingStatusNotAllowedToUpdate | 預約狀態不允許更新               |
| 404    | E3BK001 | BookingNotFound                 | 預約不存在或已被取消             |
| 500    | E9001   | SysInternalError                | 系統發生錯誤，請稍後再試         |
| 500    | E9002   | SysDatabaseError                | 資料庫操作失敗                   |

---

## 資料表

- `bookings`

---

## Service 邏輯
1. 驗證角色門市權限。
2. 驗證預約是否存在並隸屬於該門市，且預約狀態為 `COMPLETED`
3. 更新預約內容（`bookings`）。
4. 回傳最新預約資訊。

---

## 注意事項
- 預約狀態為 `COMPLETED` 才能更新完成時間。
- 目前只能更新完成時間，其他欄位無法更新。
