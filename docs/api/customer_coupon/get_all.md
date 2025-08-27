## User Story

作為顧客，我希望能夠取得我自己擁有的全部優惠券，支援查詢條件，方便於前台顯示或個人中心查閱。

---

## Endpoint

**GET** `/api/customer_coupons`

---

## 說明

- 提供顧客取得自己擁有的全部優惠券。
- 支援基本查詢條件。
- 支援分頁（limit、offset）。
- 支援排序（sort）。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Query Parameters

| 參數   | 型別   | 必填 | 預設值         | 說明                                             |
| ------ | ------ | ---- | -------------- | ------------------------------------------------ |
| isUsed | bool   | 否   |                | 是否已使用                                       |
| limit  | int    | 否   | 20             | 單頁筆數                                         |
| offset | int    | 否   | 0              | 起始筆數                                         |
| sort   | string | 否   | isUsed,validTo | 排序欄位 (可以逗號串接，有 `-` 表示 `DESC` 排序) |

### 驗證規則

| 欄位   | 必填 | 其他規則                                                                 |
| ------ | ---- | ------------------------------------------------------------------------ |
| isUsed | 否   |                                                                          |
| limit  | 否   | <li>最小值1<li>最大值100                                                 |
| offset | 否   | <li>最小值0<li>最大值1000000                                             |
| sort   | 否   | <li>可以為 createdAt, updatedAt, isUsed, validFrom, validTo (其餘會忽略) |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 3,
    "items": [
      {
        "id": "1000000001",
        "validFrom": "2025-01-01T00:00:00+08:00",
        "validTo": "", // 如果為空，表示無期限
        "isUsed": false,
        "usedAt": "", // 如果為空，表示未使用
        "createdAt": "2025-01-01T00:00:00+08:00",
        "coupon": {
          "id": "1000000001",
          "displayName": "新客優惠",
          "discountRate": 0.8,
          "discountAmount": 100,
          "isActive": true // 如果為 false，表示已失效
        },
      }
    ]
  }
}
```

---

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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                             |
| ------ | ------ | ----------------------- | -------------------------------- |
| 401    | E1002  | AuthTokenInvalid        | 無效的 accessToken，請重新登入   |
| 401    | E1003  | AuthTokenMissing        | accessToken 缺失，請重新登入     |
| 401    | E1004  | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入 |
| 401    | E1006  | AuthContextMissing      | 未找到使用者認證資訊，請重新登入 |
| 400    | E2004  | ValTypeConversionFailed | 參數類型轉換失敗                 |
| 400    | E2023  | ValFieldMinNumber       | {field} 最小值為 {param}         |
| 400    | E2026  | ValFieldMaxNumber       | {field} 最大值為 {param}         |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                   |

---

## 資料表

- `customer_coupons`
- `coupons`

---

## Service 邏輯

1. 根據 `accessToken` 取得顧客ID。
2. 查詢並回傳該顧客的全部優惠券。

---

## 注意事項

- 僅允許本人查詢。
