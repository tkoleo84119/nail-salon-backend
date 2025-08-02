## User Story

作為管理員，我希望可以查詢特定員工的門市存取權限（Store Access），以便了解其可操作的門市範圍。

---

## Endpoint

**GET** `/api/admin/staff/{staffId}/store-access`

---

## 說明

- 提供管理員查詢特定員工的門市存取權限。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN` 與 `ADMIN` 可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明        |
| ------- | ----------- |
| staffId | 員工帳號 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "storeList": [
      {
        "storeId": "8000000001",
        "name": "大安旗艦店"
      },
      {
        "storeId": "8000000003",
        "name": "信義分店"
      }
    ]
  }
}
```

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                             |
| ------ | -------- | -------------------------------- |
| 401    | E1002    | 無效的 accessToken，請重新登入   |
| 401    | E1003    | accessToken 缺失，請重新登入     |
| 401    | E1004    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010    | 權限不足，無法執行此操作         |
| 400    | E2004    | 參數類型轉換失敗                 |
| 400    | E2020    | {field} 為必填項目               |
| 404    | E3STA005 | 員工帳號不存在                   |
| 500    | E9001    | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | 資料庫操作失敗                   |

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "errors": [
    {
      "code": "E1002",
      "message": "無效的 accessToken"
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
      "message": "無權限存取此資源"
    }
  ]
}
```

#### 404 Not Found - 查無員工資料

```json
{
  "errors": [
    {
      "code": "E3STA005",
      "message": "員工帳號不存在"
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

- `staff_users`
- `staff_user_store_access`
- `stores`

---

## Service 邏輯

1. 確認 `staffId` 是否存在。
2. 查詢 `staff_user_store_access` 表取得所有該員工可操作的門市資料。
3. 回傳門市清單。

---

## 注意事項

- 若無授權任何門市，回傳空陣列。
