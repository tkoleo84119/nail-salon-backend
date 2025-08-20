## User Story

作為客戶，我想要取得我自己的資料，方便於前台顯示或個人中心查閱。

---

## Endpoint

**GET** `/api/customers/me`

---

## 說明

- 提供顧客取得自己的資料。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

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
    "isIntrovert": false,
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

| 狀態碼 | 錯誤碼 | 常數名稱               | 說明                             |
| ------ | ------ | ---------------------- | -------------------------------- |
| 401    | E1002  | AuthTokenInvalid       | 無效的 accessToken，請重新登入   |
| 401    | E1003  | AuthTokenMissing       | accessToken 缺失，請重新登入     |
| 401    | E1004  | AuthTokenFormatError   | accessToken 格式錯誤，請重新登入 |
| 401    | E1006  | AuthContextMissing     | 未找到使用者認證資訊，請重新登入 |
| 401    | E1011  | AuthCustomerFailed     | 未找到有效的顧客資訊，請重新登入 |
| 500    | E9001  | SysInternalError       | 系統發生錯誤，請稍後再試         |
| 500    | E9002  | SysDatabaseError       | 資料庫操作失敗                   |

---

## 資料表

- `customers`

---

## Service 邏輯

1. 根據 `accessToken` 取得顧客ID。
2. 查詢並回傳該顧客的完整資料。

---

## 注意事項

- 僅允許本人查詢。
