## User Story

作為一位超級管理員（`SUPER_ADMIN`）或系統管理員（`ADMIN`），我希望可以新增某位員工可操作的門市（單筆）， 以便彈性擴充該員工可管理的門市範圍。

---

## Endpoint

**POST** `/api/staff/{id}/store-access`

---

## 權限

- 僅限 `SUPER_ADMIN` 或 `ADMIN` 存取


---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Body

```json
{
  "storeId": "1"
}
```

### 驗證規則

| 欄位    | 規則     | 說明         |
| ------- | -------- | ------------ |
| storeId | <li>必填 | 欲新增的門市 |

---

## Response

### 成功 200 OK（已存在）

```json
{
  "data": {
    "staffUserId": "928374234",
    "storeList": [
      {
        "id": "1",
        "name": "台北總店"
      },
      {
        "id": "2",
        "name": "新竹巨城店"
      }
    ]
  }
}
```

### 成功 201 Created（成功新增）

```json
{
  "data": {
    "staffUserId": "928374234",
    "storeList": [
      {
        "id": "1",
        "name": "台北總店"
      },
      {
        "id": "2",
        "name": "新竹巨城店"
      }
    ]
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "storeId": "storeId為必填項目"
  }
}
```

#### 403 Forbidden - 權限不足

```json
{
  "message": "權限不足，無法執行此操作"
}
```

#### 404 Not Found - 員工不存在

```json
{
  "message": "指定的員工不存在"
}
```

#### 404 Not Found - 門市不存在

```json
{
  "message": "指定的門市不存在"
}
```

#### 500 Internal Server Error

```json
{
  "message": "系統發生錯誤，請稍後再試"
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
5. 檢查門市是否啟用中，且是否是該管理員有權限的門市
6. 查詢是否已有相同的門市權限
   - 若有：不新增，回傳 200 (全部的 store access)
   - 若無：新增一筆 `staff_user_store_access` 記錄，回傳 201 (全部的 store access)

---

## 注意事項

- 僅能操作他人帳號（不可修改自己）
- 一次僅能新增一筆 store access（非批次）
- 會回傳全部的 store access 資訊

