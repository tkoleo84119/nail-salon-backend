## User Story

作為一個客戶，我希望可以編輯自己的客戶資料，方便保持資訊正確。

---

## Endpoint

**PATCH** `/api/customers/me`

---

## 說明

- 提供顧客編輯自己的資料。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "name": "王小美",
  "phone": "0912345678",
  "birthday": "1992-02-29",
  "email": "test@test.com",
  "city": "台北市",
  "favoriteShapes": ["方形"],
  "favoriteColors": ["粉色"],
  "favoriteStyles": ["法式"],
  "isIntrovert": true,
  "customerNote": "容易指緣乾裂"
}
```

### 驗證規則

| 欄位           | 必填 | 其他規則                                                                                                    | 說明       |
| -------------- | ---- | ----------------------------------------------------------------------------------------------------------- | ---------- |
| name           | 否   | <li>不能為空字串<li>最大長度100字元                                                                         | 姓名       |
| phone          | 否   | <li>格式是09xxxxxxxx                                                                                        | 電話       |
| birthday       | 否   | <li>格式是yyyy-MM-dd                                                                                        | 生日       |
| email          | 否   | <li>email格式                                                                                               | 電子郵件   |
| city           | 否   | <li>最大長度100字元                                                                                         | 城市       |
| favoriteShapes | 否   | <li>最長20筆<li>值只能為 方形 方圓形 橢圓形 圓形 圓尖形 尖形 梯形 不一定                                    | 喜歡的指形 |
| favoriteColors | 否   | <li>最長20筆<li>值只能為 白色系 裸色系 粉色系 紅色系 橘色系 大地色系 綠色系 藍色系 紫色系 黑色系 不一定     | 喜歡色系   |
| favoriteStyles | 否   | <li>最長20筆<li>值只能為 暈染 手繪 貓眼 鏡面 可愛 法式 漸層 氣質溫柔 個性 日系 簡約 優雅 典雅 小眾 沒有固定 | 喜歡款式   |
| isIntrovert    | 否   |                                                                                                             | 是否是I人  |
| customerNote   | 否   | <li>最大長度255字元                                                                                         | 個人備註   |

- 欄位皆為選填，但至少需有一項。

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "1000000001",
    "name": "王小美",
    "phone": "0912345678",
    "birthday": "1992-02-29",
    "email": "test@test.com",
    "city": "台北市",
    "favoriteShapes": ["方形"],
    "favoriteColors": ["粉色"],
    "favoriteStyles": ["法式"],
    "isIntrovert": true,
    "customerNote": "容易指緣乾裂"
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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                                                        |
| ------ | ------ | ----------------------- | ----------------------------------------------------------- |
| 401    | E1002  | AuthTokenInvalid        | 無效的 accessToken，請重新登入                              |
| 401    | E1003  | AuthTokenMissing        | accessToken 缺失，請重新登入                                |
| 401    | E1004  | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入                            |
| 401    | E1006  | AuthContextMissing      | 未找到使用者認證資訊，請重新登入                            |
| 401    | E1011  | AuthCustomerFailed      | 未找到有效的顧客資訊，請重新登入                            |
| 400    | E2001  | ValJSONFormatError      | JSON 格式錯誤，請檢查                                       |
| 400    | E2003  | ValAllFieldsEmpty       | 至少需要提供一個欄位進行更新                                |
| 400    | E2024  | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元                       |
| 400    | E2025  | ValFieldArrayMaxLength  | {field} 最多只能有 {param} 個項目                           |
| 400    | E2027  | ValFieldInvalidEmail    | {field} 格式錯誤，請使用正確的電子郵件格式                  |
| 400    | E2030  | ValFieldOneof           | {field} 必須是 {param} 其中一個值                           |
| 400    | E2032  | ValFieldTaiwanMobile    | {field} 格式錯誤，請使用正確的台灣手機號碼格式 (0912345678) |
| 400    | E2033  | ValFieldDateFormat      | {field} 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD)         |
| 400    | E2036  | ValFieldNoBlank         | {field} 不能為空字串                                        |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試                                    |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                                              |

---

## 資料表

- `customers`

---

## Service 邏輯

1. 若有傳入 birthday，則驗證格式是否為 yyyy-MM-dd。
2. 更新 customer 資料。
3. 回傳更新後資料。

---

## 注意事項

- 僅允許本人編輯。
- 至少需要提供一個欄位進行更新。

