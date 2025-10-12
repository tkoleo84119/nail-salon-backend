## User Story

作為員工，我希望可以查詢某門市某單一產品的詳細資料，方便管理。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/products/{productId}`

---

## 說明

- 用於查詢特定產品的詳細資訊。

---

## 權限

- 需要登入才可使用。
- 所有員工都可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameters

| 參數      | 說明    |
| --------- | ------- |
| storeId   | 門市 ID |
| productId | 產品 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "9000000001",
    "name": "利德瑪",
    "brand": {
      "id": "9000000001",
      "name": "利德瑪"
    },
    "category": {
      "id": "9000000001",
      "name": "分類1"
    },
    "currentStock": 10,
    "safetyStock": 5,
    "unit": "瓶",
    "storageLocation": "櫃子B",
    "note": "左下角",
    "isActive": true,
    "createdAt": "2025-01-01T00:00:00+08:00",
    "updatedAt": "2025-01-01T00:00:00+08:00"
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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                             |
| ------ | -------- | ----------------------- | -------------------------------- |
| 401    | E1002    | AuthTokenInvalid        | 無效的 accessToken，請重新登入   |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入     |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入 |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                 |
| 400    | E3PRO003 | ProductNotBelongToStore | 產品不屬於指定的門市             |
| 404    | E3PRO002 | ProductNotFound         | 產品不存在或已被刪除             |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                   |

---

## 資料表

- `brands`
- `product_categories`
- `products`

---

## Service 邏輯

1. 查詢產品(只使用 `id` 查詢)
2. 確認產品是否屬於該門市
3. 回傳產頻詳細資訊

---

## 注意事項

- createdAt 與 updatedAt 會是標準 Iso 8601 格式。
