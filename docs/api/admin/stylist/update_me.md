## User Story

作為一位員工，我希望可以更新自己的美甲師個人資料（stylists），讓顧客可以看到我最新的專長與風格。

---

## Endpoint

**PATCH** `/api/admin/stylists/me`

---

## 說明

- 僅允許員工更新自己的美甲師資料。

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
  "stylistName": "Jane 美甲師",
  "goodAtShapes": ["方形", "橢圓形"],
  "goodAtColors": ["白色系"],
  "goodAtStyles": ["簡約", "法式"],
  "isIntrovert": true
}
```

### 驗證規則

| 欄位         | 必填 | 其他規則                                                                                           | 說明       |
| ------------ | ---- | -------------------------------------------------------------------------------------------------- | ---------- |
| stylistName  | 否   | <li>email格式                                                                                      | 員工 Email |
| goodAtShapes | 否   | <li>最多20項<li>值只能為 方形 方圓形 橢圓形 圓形 圓尖形 尖形 梯形                                  | 擅長指型   |
| goodAtColors | 否   | <li>最多20項<li>值只能為 白色系 裸色系 粉色系 紅色系 橘色系 大地色系 綠色系 藍色系 紫色系 黑色系   | 擅長色系   |
| goodAtStyles | 否   | <li>最多20項<li>值只能為 暈染 手繪 貓眼 鏡面 可愛 法式 漸層 氣質溫柔 個性 日系 簡約 優雅 典雅 小眾 | 擅長款式   |
| isIntrovert  | 否   | <li>布林值                                                                                         | 是否I人    |

- 至少需要提供一個欄位進行更新

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "18000000001",
    "staffUserId": "13984392823",
    "name": "Jane 美甲師",
    "goodAtShapes": ["方形", "橢圓形"],
    "goodAtColors": ["粉嫩系"],
    "goodAtStyles": ["簡約", "法式"],
    "isIntrovert": true,
    "createdAt": "2025-06-01T08:00:00+08:00",
    "updatedAt": "2025-06-01T08:00:00+08:00"
  }
}
```

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                                  |
| ------ | -------- | ------------------------------------- |
| 401    | E1002    | 無效的 accessToken，請重新登入        |
| 401    | E1003    | accessToken 缺失，請重新登入          |
| 401    | E1004    | accessToken 格式錯誤，請重新登入      |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入      |
| 400    | E2001    | JSON 格式錯誤，請檢查                 |
| 400    | E2003    | 至少需要提供一個欄位進行更新          |
| 400    | E2004    | 參數類型轉換失敗                      |
| 400    | E2024    | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2025    | {field} 最多只能有 {param} 個項目     |
| 400    | E2029    | {field} 必須是布林值                  |
| 400    | E2030    | {field} 必須是 {param} 其中一個值     |
| 404    | E3STA005 | 員工帳號不存在                        |
| 404    | E3STY001 | 美甲師資料不存在                      |
| 500    | E9001    | 系統發生錯誤，請稍後再試              |
| 500    | E9002    | 資料庫操作失敗                        |

#### 400 Bad Request - 驗證錯誤

```json
{
  "errors": [
    {
      "code": "E2003",
      "message": "至少需要提供一個欄位進行更新"
    }
  ]
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "errors": [
    {
      "code": "E1002",
      "message": "無效的 accessToken，請重新登入"
    }
  ]
}
```

#### 404 Not Found - 美甲師資料不存在

```json
{
  "errors": [
    {
      "code": "E3STY001",
      "message": "美甲師資料不存在"
    }
  ]
}
```

#### 500 Internal Server Error - 系統發生錯誤

```json
{
  "errors": [
    {
      "code": "E9001",
      "message": "系統發生錯誤，請稍後再試"
    }
  ]
}
```

---

## 資料表

- `stylists`
- `staff_users`

---

## Service 邏輯

1. 驗證至少一個欄位有更新。
2. 檢查 `stylists` 資料是否存在。
3. 更新 `stylists` 的指定欄位。
4. 回傳更新後的 stylist 資料。

---

## 注意事項

- 會回傳更新後的資料，不需要在呼叫一次 `GET /api/admin/staff/me` 來取得資料。
