## User Story

作為員工，我希望可以在一筆預約中新增多筆使用的產品，方便管理預約的產品使用情況。

---

## Endpoint

**POST** `/api/admin/stores/{storeId}/bookings/{bookingId}/products/bulk`

---

## 說明

- 提供後台員工在一筆預約中新增多筆使用的產品功能。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數      | 說明    |
| --------- | ------- |
| storeId   | 門市 ID |
| bookingId | 預約 ID |

### Body 範例

```json
{
  "productIds": ["9000000001", "9000000002"],
}
```

## 驗證規則

| 欄位       | 必填 | 其他規則                | 說明     |
| ---------- | ---- | ----------------------- | -------- |
| productIds | 是   | <li>最小1筆<li>最大50筆 | 產品 IDs |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "created": ["9000000001", "9000000002"]
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

| 狀態碼 | 錯誤碼   | 常數名稱                        | 說明                              |
| ------ | -------- | ------------------------------- | --------------------------------- |
| 401    | E1002    | AuthTokenInvalid                | 無效的 accessToken，請重新登入    |
| 401    | E1003    | AuthTokenMissing                | accessToken 缺失，請重新登入      |
| 401    | E1004    | AuthTokenFormatError            | accessToken 格式錯誤，請重新登入  |
| 401    | E1006    | AuthContextMissing              | 未找到使用者認證資訊，請重新登入  |
| 401    | E1011    | AuthCustomerFailed              | 未找到有效的顧客資訊，請重新登入  |
| 400    | E2001    | ValJSONFormatError              | JSON 格式錯誤，請檢查             |
| 400    | E2002    | ValPathParamMissing             | 路徑參數缺失，請檢查              |
| 400    | E2004    | ValTypeConversionFailed         | 參數類型轉換失敗                  |
| 400    | E2020    | ValFieldRequired                | {field} 為必填項目                |
| 400    | E2022    | ValFieldArrayMinLength          | {field} 至少需要 {param} 個項目   |
| 400    | E2025    | ValFieldArrayMaxLength          | {field} 最多只能有 {param} 個項目 |
| 400    | E3BK002  | BookingStatusNotAllowedToUpdate | 預約狀態不允許更新                |
| 400    | E3BK004  | BookingNotBelongToStore         | 預約不屬於指定的門市              |
| 404    | E3BK001  | BookingNotFound                 | 預約不存在或已被取消              |
| 404    | E3PRO002 | ProductNotFound                 | 產品不存在或已被刪除              |
| 500    | E9001    | SysInternalError                | 系統發生錯誤，請稍後再試          |
| 500    | E9002    | SysDatabaseError                | 資料庫操作失敗                    |

---

## 資料表

- `bookings`
- `products`
- `booking_products`

---

## Service 邏輯

1. 驗證門市權限。
2. 驗證預約是否存在。
3. 驗證預約是否屬於該門市。
4. 驗證預約狀態是否為 `COMPLETED`。
4. 驗證產品是否存在。
5. 驗證產品是否屬於該門市。
6. 查詢目前該預約使用產品，排除已經新增過的產品。
7. 建立 `booking_products` 資料。
8. 回傳資料。
