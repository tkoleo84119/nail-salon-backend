## User Story

作為一位管理員，我希望能更新產品，方便即時維護產品資訊。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}/products/{productId}`

---

## 說明

- 可更新名稱、庫存數量、安全庫存數量、單位、存放位置、備註。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` `MANAGER` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數      | 說明   |
| --------- | ------ |
| storeId   | 門市ID |
| productId | 產品ID |

### Body 範例

```json
{
  "brandId": "9000000001",
  "categoryId": "9000000001",
  "name": "凝膠",
  "currentStock": 10,
  "safetyStock": 5,
  "unit": "瓶",
  "storageLocation": "櫃子B",
  "note": "左下角",
  "isActive": true
}
```

### 驗證規則

| 欄位            | 必填 | 其他規則                            | 說明         |
| --------------- | ---- | ----------------------------------- | ------------ |
| brandId         | 否   |                                     | 品牌 ID      |
| categoryId      | 否   |                                     | 分類 ID      |
| name            | 否   | <li>不能為空字串<li>最大長度200字元 | 產品名稱     |
| currentStock    | 否   | <li>最小值0<li>最大值1000000        | 庫存數量     |
| safetyStock     | 否   | <li>最小值-1<li>最大值1000000       | 安全庫存數量 |
| unit            | 否   | <li>最大長度50字元                  | 單位         |
| storageLocation | 否   | <li>最大長度100字元                 | 存放位置     |
| note            | 否   | <li>最大長度255字元                 | 備註         |
| isActive        | 否   |                                     | 是否啟用     |

- 至少需要提供一個欄位進行更新。

---

## Response

### 成功 200 OK

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

| 狀態碼 | 錯誤碼   | 常數名稱                             | 說明                                             |
| ------ | -------- | ------------------------------------ | ------------------------------------------------ |
| 401    | E1002    | AuthTokenInvalid                     | 無效的 accessToken，請重新登入                   |
| 401    | E1003    | AuthTokenMissing                     | accessToken 缺失，請重新登入                     |
| 401    | E1004    | AuthTokenFormatError                 | accessToken 格式錯誤，請重新登入                 |
| 401    | E1005    | AuthStaffFailed                      | 未找到有效的員工資訊，請重新登入                 |
| 401    | E1006    | AuthContextMissing                   | 未找到使用者認證資訊，請重新登入                 |
| 403    | E1010    | AuthPermissionDenied                 | 權限不足，無法執行此操作                         |
| 400    | E2001    | ValJsonFormat                        | JSON 格式錯誤，請檢查                            |
| 400    | E2002    | ValPathParamMissing                  | 路徑參數缺失，請檢查                             |
| 400    | E2003    | ValAllFieldsEmpty                    | 至少需要提供一個欄位進行更新                     |
| 400    | E2004    | ValTypeConversionFailed              | 參數類型轉換失敗                                 |
| 400    | E2023    | ValFieldMinNumber                    | {field} 最小值為 {param}                         |
| 400    | E2024    | ValFieldStringMaxLength              | {field} 長度最多只能有 {param} 個字元            |
| 400    | E2026    | ValFieldMaxNumber                    | {field} 最大值為 {param}                         |
| 400    | E2036    | ValFieldNoBlank                      | {field} 不能為空字串                             |
| 400    | E3PRO003 | ProductNotBelongToStore              | 產品不屬於指定的門市                             |
| 404    | E3PRO002 | ProductNotFound                      | 產品不存在或已被刪除                             |
| 404    | E3BRN002 | BrandNotFound                        | 品牌不存在或已被刪除                             |
| 404    | E3PC002  | CategoryNotFound                     | 分類不存在或已被刪除                             |
| 409    | E3PRO001 | ProductNameBrandAlreadyExistsInStore | 產品名稱和品牌組合已存在於同門市，請使用其他組合 |
| 500    | E9001    | SysInternalError                     | 系統發生錯誤，請稍後再試                         |
| 500    | E9002    | SysDatabaseError                     | 資料庫操作失敗                                   |

---

## 資料表

- `brands`
- `product_categories`
- `products`

---

## Service 邏輯

1. 驗證門市存取權限。
2. 確認產品是否存在。
   - 確認產品是否屬於該門市。
3. 若有更新 `brandId`、`categoryId`，則驗證 `brandId`、`categoryId` 是否存在。
4. 若有更新 `name` 或 `brandId`，則驗證 `storeId`、`name`、`brandId` 是否唯一 (不包含自己)。
5. 更新 `products` 資料。
6. 回傳更新結果。

---

## 注意事項

- 同一門市下，產品名稱 + 品牌 不可重複。
