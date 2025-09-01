package sqlx

import (
	"github.com/jmoiron/sqlx"
)

// Repositories consolidates all SQLX repositories into a single interface
type Repositories struct {
	Booking        *BookingRepository
	BookingDetail  *BookingDetailRepository
	Brand          *BrandRepository
	Customer       *CustomerRepository
	Coupon         *CouponRepository
	CustomerCoupon *CustomerCouponRepository
	ProductCategory *ProductCategoryRepository
	Schedule       *ScheduleRepository
	Service        *ServiceRepository
	Staff          *StaffUserRepository
	Store          *StoreRepository
	Stylist        *StylistRepository
	TimeSlot       *TimeSlotRepository
	Template       *TimeSlotTemplateRepository
}

// NewRepositories creates a new instance of Repositories with all repository instances
func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Booking:        NewBookingRepository(db),
		BookingDetail:  NewBookingDetailRepository(db),
		Brand:          NewBrandRepository(db),
		Customer:       NewCustomerRepository(db),
		Coupon:         NewCouponRepository(db),
		CustomerCoupon: NewCustomerCouponRepository(db),
		ProductCategory: NewProductCategoryRepository(db),
		Schedule:       NewScheduleRepository(db),
		Service:        NewServiceRepository(db),
		Staff:          NewStaffUserRepository(db),
		Store:          NewStoreRepository(db),
		Stylist:        NewStylistRepository(db),
		TimeSlot:       NewTimeSlotRepository(db),
		Template:       NewTimeSlotTemplateRepository(db),
	}
}
