## User Story

作為一位超級管理員（`SUPER_ADMIN`）或系統管理員（`ADMIN`），我希望可以刪除某位員工可操作的門市（多筆），以便控管其實際管理範圍。

---

## Endpoint

**DELETE** `/api/staff/{id}/store-access`

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
  "storeIds": ["1", "3"]
}
```

### 驗證規則

| 欄位     | 規則                         | 說明                 |
| -------- | ---------------------------- | -------------------- |
| storeIds | <li>必填<li>陣列<li>不可為空 | 欲移除的門市 ID 清單 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "staffUserId": "928374234",
    "storeList": [
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
    "storeIds": "門市清單最小值為1"
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
6. 執行刪除動作（從 `staff_user_store_access` 中移除多筆）
7. 回傳目前該員工剩餘可操作的門市清單

---

## 注意事項

- 僅能操作他人帳號（不可移除自己）
- `SUPER_ADMIN` 不能被更動其門市權限
- 不會因門市不存在於原 access 而報錯（可忽略）

