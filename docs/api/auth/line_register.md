## User Story

作為一個客戶，當我要完成 LINE 註冊時，系統會將我的資訊存入後，回傳 access token、refresh token 等。

---

## Endpoint

**POST** `/api/auth/line/register`

---

## 說明

- 用戶在 LINE 登入後，若尚未註冊，需呼叫本 API 完成註冊流程。
- 完成註冊後，自動發 access token 與 refresh token。
- 需帶入從 LINE profile 取得的必要資訊與額外註冊欄位。

---

## 權限

- 不須預先認證

---

## Request

### Header

- Content-Type: application/json

### Body 範例

```json
{
  "idToken": "eyJraWQiOiJ...",
  "name": "小美",
  "phone": "09xxxxxxxx",
  "birthday": "1990-01-01",
  "city": "台北市",
  "favoriteShapes": ["圓形", "方形"],
  "favoriteColors": ["黑色", "白色"],
  "favoriteStyles": ["自然", "韓式"],
  "isIntrovert": true,
  "referralSource": ["朋友介紹", "網路廣告"],
  "referrer": "1000000001",
  "customerNote": "這是客戶的備註",
}
```

### 驗證規則

| 欄位           | 必填 | 其他規則                                                                                                    | 說明             |
| -------------- | ---- | ----------------------------------------------------------------------------------------------------------- | ---------------- |
| idToken        | 是   | <li>最大長度500字元                                                                                         | LINE idToken     |
| name           | 是   | <li>最大長度100字元                                                                                         | 姓名             |
| phone          | 是   | <li>格式是09xxxxxxxx                                                                                        | 電話             |
| birthday       | 是   | <li>格式是yyyy-MM-dd                                                                                        | 生日             |
| city           | 否   | <li>最大長度100字元                                                                                         | 城市             |
| favoriteShapes | 否   | <li>最多20項<li>值只能為 方形 方圓形 橢圓形 圓形 圓尖形 尖形 梯形 不一定                                    | 喜歡的指形       |
| favoriteColors | 否   | <li>最多20項<li>值只能為 白色系 裸色系 粉色系 紅色系 橘色系 大地色系 綠色系 藍色系 紫色系 黑色系  不一定    | 喜歡的色系       |
| favoriteStyles | 否   | <li>最多20項<li>值只能為 暈染 手繪 貓眼 鏡面 可愛 法式 漸層 氣質溫柔 個性 日系 簡約 優雅 典雅 小眾 沒有固定 | 喜歡的款式       |
| isIntrovert    | 否   | <li>布林值                                                                                                  | 是否是I人        |
| referralSource | 否   | <li>最多20項<li>值只能為 Facebook Instagram Threads Dcard Google 親友介紹                                   | 推薦來源         |
| referrer       | 否   | <li>最大長度100字元                                                                                         | 推薦人           |
| customerNote   | 否   | <li>最大長度255字元                                                                                         | 使用者自己的備註 |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "accessToken": "...",
    "refreshToken": "...",
    "expiresIn": 3600
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
| 401    | E1007  | AuthLineTokenInvalid    | Line idToken 驗證失敗，請重新登入                           |
| 401    | E1008  | AuthLineTokenExpired    | Line idToken 已過期，請重新登入                             |
| 400    | E2001  | ValJsonFormat           | JSON 格式錯誤，請檢查                                       |
| 400    | E2004  | ValTypeConversionFailed | 參數類型轉換失敗                                            |
| 400    | E2020  | ValFieldRequired        | {field} 為必填項目                                          |
| 400    | E2024  | ValFieldMaxLength       | {field} 長度最多只能有 {param} 個字元                       |
| 400    | E2025  | ValFieldArrayMaxLength  | {field} 最多只能有 {param} 個項目                           |
| 400    | E2029  | ValFieldBoolean         | {field} 必須是布林值                                        |
| 400    | E2030  | ValFieldOneOf           | {field} 必須是 {param} 其中一個值                           |
| 400    | E2032  | ValFieldTaiwanMobile    | {field} 格式錯誤，請使用正確的台灣手機號碼格式 (0912345678) |
| 400    | E2033  | ValFieldDateFormat      | {field} 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD)         |
| 409    | E3C003 | CustomerAlreadyExists   | 客戶已存在                                                  |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試                                    |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                                              |

---

## 資料表

- `customers`
- `customer_tokens`

---

## Service 邏輯

1. 呼叫 LINE 驗證 `idToken` 合法性，取得 `providerUid`。
2. 驗證該 `providerUid` 是否已註冊（重複則 409）。
3. 建立 `customers` 資料與對應 `customer_tokens`。
4. 產生 `access token`、`refresh token`。
5. 回傳 `access token`、`refresh token`。

---

## 注意事項

- `level` 預設為 `NORMAL`。
