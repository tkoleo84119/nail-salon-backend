package sqlx

import (
	"github.com/jmoiron/sqlx"
)

// Repositories consolidates all SQLX repositories into a single interface
type Repositories struct {
	Account         *AccountRepository
	Booking         *BookingRepository
	BookingDetail   *BookingDetailRepository
	BookingProduct  *BookingProductRepository
	Brand           *BrandRepository
	Customer        *CustomerRepository
	Coupon          *CouponRepository
	CustomerCoupon  *CustomerCouponRepository
	Expense         *ExpenseRepository
	ExpenseItem     *ExpenseItemRepository
	Product         *ProductRepository
	ProductCategory *ProductCategoryRepository
	Schedule        *ScheduleRepository
	Service         *ServiceRepository
	Staff           *StaffUserRepository
	StockUsage      *StockUsageRepository
	Store           *StoreRepository
	Stylist         *StylistRepository
	Supplier        *SupplierRepository
	TimeSlot        *TimeSlotRepository
	Template        *TimeSlotTemplateRepository
}

// NewRepositories creates a new instance of Repositories with all repository instances
func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Account:         NewAccountRepository(db),
		Booking:         NewBookingRepository(db),
		BookingDetail:   NewBookingDetailRepository(db),
		BookingProduct:  NewBookingProductRepository(db),
		Brand:           NewBrandRepository(db),
		Customer:        NewCustomerRepository(db),
		Coupon:          NewCouponRepository(db),
		CustomerCoupon:  NewCustomerCouponRepository(db),
		Expense:         NewExpenseRepository(db),
		ExpenseItem:     NewExpenseItemRepository(db),
		Product:         NewProductRepository(db),
		ProductCategory: NewProductCategoryRepository(db),
		Schedule:        NewScheduleRepository(db),
		Service:         NewServiceRepository(db),
		Staff:           NewStaffUserRepository(db),
		StockUsage:      NewStockUsageRepository(db),
		Store:           NewStoreRepository(db),
		Stylist:         NewStylistRepository(db),
		Supplier:        NewSupplierRepository(db),
		TimeSlot:        NewTimeSlotRepository(db),
		Template:        NewTimeSlotTemplateRepository(db),
	}
}
