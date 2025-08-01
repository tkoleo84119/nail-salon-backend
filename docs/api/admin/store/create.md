## User Story

作為一位管理員，我希望能新增門市（store），方便多據點管理。

---

## Endpoint

**POST** `/api/admin/stores`

---

## 說明

- 提供後台管理員新增門市功能。
- `ADMIN` 新增門市時會自動取得該門市權限。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "name": "大安旗艦店",
  "address": "台北市大安區復興南路一段100號",
  "phone": "02-12345678"
}
```

### 驗證規則

| 欄位    | 規則                                                           | 說明     |
| ------- | -------------------------------------------------------------- | -------- |
| name    | <li>必填<li>長度小於100                                        | 門市名稱 |
| address | <li>選填<li>長度小於255                                        | 門市地址 |
| phone   | <li>選填<li>長度小於20<li>格式必須為台灣市話號碼 (02-12345678) | 電話     |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "8000000001",
    "name": "大安旗艦店",
    "address": "台北市大安區復興南路一段100號",
    "phone": "02-12345678",
    "isActive": true,
  }
}
```

### 錯誤處理

| 狀態碼 | 錯誤碼   | 說明                                                         |
| ------ | -------- | ------------------------------------------------------------ |
| 401    | E1002    | 無效的 accessToken，請重新登入                               |
| 401    | E1003    | accessToken 缺失，請重新登入                                 |
| 401    | E1004    | accessToken 格式錯誤，請重新登入                             |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入                             |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入                             |
| 403    | E1010    | 權限不足，無法執行此操作                                     |
| 400    | E2001    | JSON 格式錯誤，請檢查                                        |
| 400    | E2020    | {field} 為必填項目                                           |
| 400    | E2024    | {field} 長度最多只能有 {param} 個字元                        |
| 400    | E2031    | {field} 格式錯誤，請使用正確的台灣電話號碼格式 (0X-XXXXXXXX) |
| 409    | E3STO003 | 門市已存在，請創建其他門市                                   |
| 500    | E9001    | 系統發生錯誤，請稍後再試                                     |
| 500    | E9002    | 資料庫操作失敗                                               |

#### 400 Bad Request - 輸入驗證失敗

```json
{
  "errors": [
    {
      "code": "E2020",
      "message": "name 欄位為必填項目",
      "field": "name"
    }
  ]
}
```

#### 401 Unauthorized - 認證失敗

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

#### 403 Forbidden - 權限不足

```json
{
  "errors": [
    {
      "code": "E1010",
      "message": "權限不足，無法執行此操作"
    }
  ]
}
```

#### 500 Internal Server Error - 系統錯誤

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

- `stores`
- `staff_user_store_access`

---

## Service 邏輯

1. 再次驗證 `role` 是否為 `SUPER_ADMIN` 或 `ADMIN`。
2. 驗證 `name` 是否存在。
3. 建立 `stores` 資料。
4. 若為 `ADMIN` 則於 `staff_user_store_access` 關聯該門市。
5. 回傳新增結果。

---

## 注意事項

- 門市名稱不可重複。
- `ADMIN` 新增門市時會自動取得該門市權限。
