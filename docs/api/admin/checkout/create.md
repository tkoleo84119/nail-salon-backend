## User Story

作為一位員工，我希望能針對預約進行結帳，方便維護結帳資訊。

---

## Endpoint

**POST** `/api/admin/stores/:storeId/bookings/:bookingId/checkouts`

---

## 說明

- 提供員工針對預約進行結帳功能。

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
  "paymentMethod": "cash",
  "couponId": "1234567890",
  "paidAmount": 1000,
  "bookingDetails": [
    {
      "id": "1234567890",
      "price": 500,
      "useCoupon": true,
    }
  ]
}
```

### 驗證規則

| 欄位                     | 必填 | 其他規則                          | 說明           |
| ------------------------ | ---- | --------------------------------- | -------------- |
| paymentMethod            | 是   | <li>值可以為 `cash` `linePay`     | 付款方式       |
| customerCouponId         | 否   |                                   | 客戶優惠券ID   |
| paidAmount               | 是   | <li>最小值為 0<li>最大值為1000000 | 實際付款金額   |
| bookingDetails           | 否   | <li>最少1筆<li>最多10筆           | 預約明細       |
| bookingDetails.id        | 是   |                                   | 預約明細ID     |
| bookingDetails.price     | 是   |                                   | 預約明細價格   |
| bookingDetails.useCoupon | 是   |                                   | 是否使用優惠券 |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "9000000001"
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

| 狀態碼 | 錯誤碼    | 常數名稱                          | 說明                              |
| ------ | --------- | --------------------------------- | --------------------------------- |
| 401    | E1002     | AuthTokenInvalid                  | 無效的 accessToken，請重新登入    |
| 401    | E1003     | AuthTokenMissing                  | accessToken 缺失，請重新登入      |
| 401    | E1004     | AuthTokenFormatError              | accessToken 格式錯誤，請重新登入  |
| 401    | E1005     | AuthStaffFailed                   | 未找到有效的員工資訊，請重新登入  |
| 401    | E1006     | AuthContextMissing                | 未找到使用者認證資訊，請重新登入  |
| 403    | E1010     | AuthPermissionDenied              | 權限不足，無法執行此操作          |
| 400    | E2001     | ValJsonFormat                     | JSON 格式錯誤，請檢查             |
| 400    | E2002     | ValPathParamMissing               | 路徑參數缺失，請檢查              |
| 400    | E2004     | ValTypeConversionFailed           | 參數類型轉換失敗                  |
| 400    | E2020     | ValFieldRequired                  | {field} 為必填項目                |
| 400    | E2022     | ValFieldArrayMinLength            | {field} 至少需要 {param} 個項目   |
| 400    | E2023     | ValFieldMinNumber                 | {field} 最小值為 {param}          |
| 400    | E2025     | ValFieldArrayMaxLength            | {field} 最多只能有 {param} 個項目 |
| 400    | E2026     | ValFieldMaxNumber                 | {field} 最大值為 {param}          |
| 400    | E2030     | ValFieldOneof                     | {field} 必須是 {param} 其中一個值 |
| 400    | E3BK004   | BookingNotBelongToStore           | 預約不屬於指定的門市              |
| 400    | E3BK008   | BookingStatusNotCheckout          | 預約狀態不允許結帳                |
| 400    | E3CCOU001 | CustomerCouponNotBelongToCustomer | 客戶優惠券不屬於指定的顧客        |
| 400    | E3CCOU002 | CustomerCouponAlreadyUsed         | 客戶優惠券已使用                  |
| 400    | E3CCOU003 | CustomerCouponExpired             | 客戶優惠券已過期                  |
| 400    | E3COU001  | CouponNotActive                   | 優惠券未啟用                      |
| 404    | E3BKD001  | BookingDetailNotFound             | 預約明細不存在或已被刪除          |
| 500    | E9001     | SysInternalError                  | 系統發生錯誤，請稍後再試          |
| 500    | E9002     | SysDatabaseError                  | 資料庫操作失敗                    |

---

## 資料表

- `checkouts`
- `bookings`
- `booking_details`
- `coupons`
- `customers`

---

## Service 邏輯

1. 確認該使用者是否有權限操作該門市。
2. 確認 `booking_id` 是否存在。
3. 確認 `booking_id` 狀態是否是 `SCHEDULED`。
4. 確認 `booking_id` 是否屬於該門市。
5. 若有傳入 `customerCouponId`，則確認 `customerCouponId` 是否存在。
   - 確認 `customerCouponId` 是否跟 `bookings` 的 `customer_id` 相同。
   - 確認仍未被使用。
   - 確認優惠券是否啟用。
   - 確認優惠券是否過期。
6. 取出所有 `booking_details` 資料。
7. 比對傳入的 `booking_details` 是否存在於 `booking_details` 中，並且根據有無使用優惠券，更新 `booking_details` 的 `price`、`discount_rate`、`discount_amount`。
8. 準備 `checkouts` 資料。
9. 建立 `checkouts` 資料。
10. 更新 `booking_details` 資料。
11. 更新 `bookings` 狀態為 `COMPLETED`。
12. 若有使用優惠券，則更新 `customer_coupons` 為已使用。
13. 回傳新增結果。

---

## 注意事項

- `paymentMethod` 未來可能會再擴充。
