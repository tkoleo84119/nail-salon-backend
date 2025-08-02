## User Story

作為一位管理員，我希望可以新增某位員工可操作的門市（單筆）， 以便彈性擴充該員工可管理的門市範圍。

---

## Endpoint

**POST** `/api/admin/staff/{staffId}/store-access`

---

## 說明

- 提供管理員新增特定員工的門市存取權限。

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

### Body 範例

```json
{
  "storeId": "1"
}
```

### 驗證規則

| 欄位    | 必填 | 其他規則 | 說明         |
| ------- | ---- | -------- | ------------ |
| storeId | 是   |          | 欲新增的門市 |

---

## Response

### 成功 200 OK（已存在）

```json
{
  "data": {
    "storeList": [
      {
        "id": "1",
        "name": "台北總店"
      }
    ]
  }
}
```

### 成功 201 Created（成功新增）

```json
{
  "data": {
    "storeList": [
      {
        "id": "1",
        "name": "台北總店"
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
| 400    | E2002    | 路徑參數缺失，請檢查             |
| 400    | E2004    | 參數類型轉換失敗                 |
| 400    | E2020    | {field} 為必填項目               |
| 400    | E3STA004 | 不可更新自己的帳號               |
| 404    | E3STA005 | 員工帳號不存在                   |
| 404    | E3STO002 | 門市不存在或已被刪除             |
| 500    | E9001    | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | 資料庫操作失敗                   |

#### 400 Bad Request - 驗證錯誤

```json
{
  "errors": [
    {
        "code": "E2002",
        "message": "路徑參數缺失，請檢查",
        "field": "staffId"
    }
  ]
  }
```

#### 401 Unauthorized - 無效的 accessToken

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

#### 404 Not Found - 員工不存在

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

#### 404 Not Found - 門市不存在

```json
{
  "errors": [
    {
        "code": "E3STO002",
        "message": "指定的門市不存在或已被刪除"
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


1. 驗證目標員工是否存在
2. 不能新增自己的 store access
3. 目標員工不能為 `SUPER_ADMIN`
4. 驗證門市是否存在
5. 檢查該門市是否為該管理員有權限的門市
6. 查詢是否已有相同的門市權限
   - 若有：不新增，回傳 200 (全部的 store access)
   - 若無：新增一筆 `staff_user_store_access` 記錄，回傳 201 (全部的 store access)

---

## 注意事項

- 一次僅能新增一筆 store access
- 會回傳全部的 store access 資訊
