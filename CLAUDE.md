# CLAUDE.md

## 專案概覽
- **專案名稱**: 美甲沙龍管理系統後端 API
- **主要語言**: Go 1.24.5
- **核心框架**: Gin HTTP Web Framework
- **資料庫**: PostgreSQL (pgx/v5 + sqlx)
- **認證系統**: JWT + LINE OAuth 2.0
- **主要功能**: 預約管理、員工管理、店鋪管理、顧客管理

## 技術架構
```
nail-salon-backend/
├── cmd/server/         # 應用程式進入點
├── internal/           # 核心業務程式碼
│   ├── app/           # 應用容器與路由
│   ├── handler/       # HTTP 處理層 (Controller)
│   ├── service/       # 業務邏輯層
│   ├── model/         # 資料傳輸物件 (DTO)
│   ├── repository/    # 資料存取層
│   ├── middleware/    # HTTP 中介軟體
│   ├── infra/         # 基礎設施層
│   ├── errors/        # 錯誤處理
│   └── utils/         # 工具函式
├── migration/         # 資料庫遷移
├── scripts/          # 工具腳本
└── docs/             # API 文件
```

## 開發指令
```bash
# 應用程式執行
make run                # 啟動伺服器
go run cmd/server/main.go

# 測試相關
make test              # 執行所有測試
go test ./...

# 資料庫操作
make migrate-up        # 執行資料庫遷移
make migrate-down      # 回滾資料庫遷移
make seed-test         # 植入測試資料
```

## API 開發模式

### 標準 3 層架構模式
所有 API 都遵循：Handler → Service → Repository/SQLC

#### API 操作類型開發指南

**1. CREATE 操作** (@docs/api/*/create.md)
- 處理 POST 請求，新增資源
- 驗證請求資料與轉換 ID 類型
- 權限檢查與業務邏輯驗證
- 回傳 201 Created

**2. GET_ALL 操作** (@docs/api/*/get_all.md)
- 處理 GET 請求，查詢資源列表
- 分頁參數處理 (limit, offset)
- 排序與篩選條件
- 回傳 200 OK

**3. GET 單一資源** (@docs/api/*/get.md)
- 處理 GET 請求，查詢單一資源
- 路徑參數驗證
- 權限檢查（僅能查看自己的資料）
- 回傳 200 OK

**4. UPDATE 操作** (@docs/api/*/update.md)
- 處理 PUT/PATCH 請求，更新資源
- 部分更新支援 (HasUpdates 檢查)
- 業務邏輯驗證
- 回傳 200 OK

**5. DELETE 操作**
- 處理 DELETE 請求，刪除資源
- 軟刪除或硬刪除
- 級聯刪除檢查
- 回傳 204 No Content

**6. GET_ME / UPDATE_ME 操作** (@docs/api/*/get_me.md)
- 基於認證上下文的個人資料操作
- 自動從 JWT 中提取使用者 ID
- 無需額外權限檢查

**7. BULK 操作** (@docs/api/*/create_bulk.md, @docs/api/*/delete_bulk.md)
- 批次處理多筆資料
- 交易包裝確保原子性
- 批次驗證與錯誤處理

### 檔案命名規則
```
internal/
├── handler/{context}/{domain}/
│   ├── create.go, get.go, get_all.go
│   ├── update.go, delete.go
│   ├── get_me.go, update_me.go
│   └── {operation}_bulk.go
├── service/{context}/{domain}/
│   ├── interface.go      # 介面定義
│   └── [同上操作檔案]
└── model/{context}/{domain}/
    └── [同上操作檔案]
```

### 層級職責劃分

#### Handler 層職責
- **請求驗證**: 綁定與驗證 JSON/Query 參數
- **認證檢查**: 從 JWT 中提取 Customer/Staff 上下文
- **參數轉換**: 字串 ID 轉 int64，日期字串轉 Time
- **回應處理**: HTTP 狀態碼與統一格式回應
- **錯誤轉換**: Service 錯誤轉 HTTP 錯誤回應

#### Service 層職責
- **業務邏輯**: 實作核心業務規則與驗證
- **權限控制**: 角色權限檢查與資料存取授權 (與邏輯相關的權限檢查，如：該角色是否能夠存取某 store 的資料)
- **資料協調**: 協調多個資料來源或外部服務
- **交易管理**: 管理資料庫交易與一致性
- **資料轉換**: DB 物件轉 Response 物件

```go
// Handler 結構
type {Operation} struct {
    service {domain}Service.{Operation}Interface
}

// 標準處理流程
func (h *{Operation}) {Operation}(c *gin.Context) {
    // 1. 參數解析與驗證
    // 2. 認證上下文提取
    // 3. ID 類型轉換
    // 4. 服務層呼叫
    // 5. 回應處理
}
```

## 認證與授權
- **顧客認證**: LINE OAuth → `GetCustomerFromContext(c)`
- **員工認證**: 帳密登入 → `GetStaffFromContext(c)`
- **角色層級**: `SUPER_ADMIN` > `ADMIN` > `MANAGER` > `STYLIST`

## 錯誤處理
```go
// 統一錯誤回應
errorCodes.RespondWithServiceError(c, err)
errorCodes.RespondWithValidationErrors(c, validationErrors)
errorCodes.AbortWithError(c, errorCodes.{ErrorType}, details)
```

## 資料庫使用策略

### SQLC 使用時機
- **簡單 CRUD**: 標準新增、查詢、更新、刪除
- **靜態查詢**: 查詢條件固定，可預先定義
- **型別安全**: 自動生成型別安全的查詢函式

### SQLX 使用時機
- **動態查詢**: WHERE 條件根據參數動態組合 (如: `get_all`)
- **部分更新**: 僅更新有變化的欄位 (如: `update`)

### SQLX 開發模式
```go
// 1. Repository 結構定義
type {Domain}Repository struct {
    db *sqlx.DB
}

// 2. GET_ALL 動態查詢模式
type GetAll{Domain}sByFilterParams struct {
    Name     *string  // 使用指標支援 nil 檢查
    IsActive *bool
    Limit    *int
    Offset   *int
    Sort     *[]string
}

func (r *{Domain}Repository) GetAll{Domain}sByFilter(ctx context.Context, params GetAll{Domain}sByFilterParams) (int, []GetAll{Domain}sByFilterItem, error) {
    // 動態組裝 WHERE 條件
    whereConditions := []string{}
    args := []interface{}{}

    if params.Name != nil && *params.Name != "" {
        whereConditions = append(whereConditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
        args = append(args, "%"+*params.Name+"%")
    }

    // 先計數後查詢的效能模式
    // 分頁與排序處理
    // 執行查詢
}

// 3. UPDATE 部分欄位更新模式
type Update{Domain}Params struct {
    Name     *string  // 所有欄位都使用指標
    IsActive *bool
}

func (r *{Domain}Repository) Update{Domain}(ctx context.Context, id int64, params Update{Domain}Params) (Update{Domain}Response, error) {
    setParts := []string{"updated_at = NOW()"}
    args := []interface{}{}

    // 動態組裝 SET 條件
    if params.Name != nil && *params.Name != "" {
        setParts = append(setParts, fmt.Sprintf("name = $%d", len(args)+1))
        args = append(args, *params.Name)
    }

    // 檢查是否有欄位需要更新
    if len(setParts) == 1 {
        return Update{Domain}Response{}, fmt.Errorf("no fields to update")
    }
}

// 4. 交易操作模式
func (r *{Domain}Repository) {Operation}Tx(ctx context.Context, tx *sqlx.Tx, params {Operation}TxParams) error {
    // 使用傳入的 transaction 執行操作
}
```

### 交易（Transaction）使用指引

#### 使用 PGX Pool 交易時機
- **跨多個 SQLC 操作**: 多個 SQLC 查詢需要原子性
- **BULK 操作**: 批次新增、更新、刪除
- **複雜業務流程**: 涉及多個資料表的複雜操作
- **需要回滾**: 操作失敗時需要回滾所有變更

#### 使用 SQLX 交易時機
- **條件式更新**: 基於複雜條件的資料更新操作 (如: `update`)

## 工具函式
```go
// ID 處理
utils.GenerateID()              // 產生 Snowflake ID
utils.ParseID(stringID)         // 字串轉 int64
utils.FormatID(int64ID)         // int64 轉字串

// 時間處理
utils.TimeToPgTimestamptz(time.Now())
utils.PgTimestamptzToTimeString(pgTime)

// 分頁處理
utils.SetDefaultValuesOfPagination(limit, offset, 20, 0)
```

## 開發流程
1. 確認 API 類型與上下文 (admin/customer)
2. 建立 Model 結構 (@docs/api/)
3. 定義 Service 介面
4. 實作 Service 邏輯
5. 建立 Handler 處理
6. 註冊容器與路由
7. 執行測試驗證

## 環境變數
```bash
DB_URL=postgresql://...
JWT_SECRET=your-secret-key
LINE_CHANNEL_ID=your-line-channel-id
LINE_CHANNEL_SECRET=your-line-channel-secret
PORT=8080
```

## 詳細文件引用
- **API 規格**: @docs/api/{domain}/{operation}.md
- **資料庫設計**: @docs/db/database.dbml
- **中介軟體**: @docs/middleware/*.md
- **SQLX 開發指南**: @docs/repository/sqlx_guide.md

---
**重要提醒**:
- 所有 API 都必須遵循統一的 3 層架構模式
- 確保型別安全的 ID 轉換處理
- 實作適當的權限檢查與錯誤處理
- 使用 SQLC 進行資料庫操作的型別安全保證