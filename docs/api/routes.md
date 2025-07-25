# API Routes Documentation

This document lists all API routes in the nail salon backend system, organized by public/customer routes and admin routes.

## Public/Customer Routes

### Authentication
| Method | Endpoint                  | Description                | Status        |
| ------ | ------------------------- | -------------------------- | ------------- |
| POST   | `/api/auth/line/login`    | Customer LINE login        | ✅ Implemented |
| POST   | `/api/auth/line/register` | Customer LINE registration | ✅ Implemented |

### Customer Profile
| Method | Endpoint            | Description       | Status        |
| ------ | ------------------- | ----------------- | ------------- |
| GET    | `/api/customers/me` | Get my profile    | ✅ Implemented |
| PATCH  | `/api/customers/me` | Update my profile | ✅ Implemented |

### Booking Management
| Method | Endpoint                          | Description         | Status        |
| ------ | --------------------------------- | ------------------- | ------------- |
| POST   | `/api/bookings`                   | Create booking      | ✅ Implemented |
| GET    | `/api/bookings`                   | List my bookings    | ✅ Implemented |
| GET    | `/api/bookings/:bookingId`        | Get booking details | 🔄 TODO        |
| PATCH  | `/api/bookings/:bookingId`        | Update my booking   | ✅ Implemented |
| PATCH  | `/api/bookings/:bookingId/cancel` | Cancel my booking   | ✅ Implemented |

### Browse Services (Read-only)
| Method | Endpoint                   | Description         | Status |
| ------ | -------------------------- | ------------------- | ------ |
| GET    | `/api/services/:serviceId` | Get service details | 🔄 TODO |

### Browse Stores (Read-only)
| Method | Endpoint                        | Description         | Status |
| ------ | ------------------------------- | ------------------- | ------ |
| GET    | `/api/stores`                   | List stores         | 🔄 TODO |
| GET    | `/api/stores/:storeId`          | Get store details   | 🔄 TODO |
| GET    | `/api/stores/:storeId/stylists` | List store stylists | 🔄 TODO |
| GET    | `/api/stores/:storeId/services` | List store services | 🔄 TODO |

### Browse Schedules & Time Slots (Read-only)
| Method | Endpoint                                | Description               | Status |
| ------ | --------------------------------------- | ------------------------- | ------ |
| GET    | `/api/stores/:storeId/schedules`        | List store schedules      | 🔄 TODO |
| GET    | `/api/schedules/:scheduleId/time-slots` | List available time slots | 🔄 TODO |

## Admin Routes

### Authentication
| Method | Endpoint                | Description | Status        |
| ------ | ----------------------- | ----------- | ------------- |
| POST   | `/api/admin/auth/login` | Staff login | ✅ Implemented |

### Staff Management
| Method | Endpoint                    | Description       | Status        |
| ------ | --------------------------- | ----------------- | ------------- |
| GET    | `/api/admin/staff`          | List all staff    | 🔄 TODO        |
| POST   | `/api/admin/staff`          | Create staff      | ✅ Implemented |
| GET    | `/api/admin/staff/me`       | Get my profile    | 🔄 TODO        |
| PATCH  | `/api/admin/staff/me`       | Update my profile | ✅ Implemented |
| GET    | `/api/admin/staff/:staffId` | Get staff details | 🔄 TODO        |
| PATCH  | `/api/admin/staff/:staffId` | Update staff      | ✅ Implemented |

### Staff Store Access
| Method | Endpoint                                      | Description                | Status        |
| ------ | --------------------------------------------- | -------------------------- | ------------- |
| GET    | `/api/admin/staff/:staffId/store-access`      | List staff store access    | 🔄 TODO        |
| POST   | `/api/admin/staff/:staffId/store-access`      | Grant store access         | ✅ Implemented |
| DELETE | `/api/admin/staff/:staffId/store-access/bulk` | Revoke store access (bulk) | ✅ Implemented |

### Store Management
| Method | Endpoint                     | Description       | Status        |
| ------ | ---------------------------- | ----------------- | ------------- |
| GET    | `/api/admin/stores`          | List all stores   | 🔄 TODO        |
| POST   | `/api/admin/stores`          | Create store      | ✅ Implemented |
| GET    | `/api/admin/stores/:storeId` | Get store details | 🔄 TODO        |
| PATCH  | `/api/admin/stores/:storeId` | Update store      | ✅ Implemented |

### Service Management
| Method | Endpoint                         | Description         | Status        |
| ------ | -------------------------------- | ------------------- | ------------- |
| GET    | `/api/admin/services`            | List all services   | 🔄 TODO        |
| POST   | `/api/admin/services`            | Create service      | ✅ Implemented |
| GET    | `/api/admin/services/:serviceId` | Get service details | 🔄 TODO        |
| PATCH  | `/api/admin/services/:serviceId` | Update service      | ✅ Implemented |
| DELETE | `/api/admin/services/:serviceId` | Deactivate service  | 🔄 TODO        |

### Stylist Management
| Method | Endpoint                         | Description               | Status        |
| ------ | -------------------------------- | ------------------------- | ------------- |
| GET    | `/api/admin/stylists`            | List all stylists         | 🔄 TODO        |
| GET    | `/api/admin/stylists/me`         | Get my stylist profile    | 🔄 TODO        |
| POST   | `/api/admin/stylists/me`         | Create my stylist profile | ✅ Implemented |
| PATCH  | `/api/admin/stylists/me`         | Update my stylist profile | ✅ Implemented |
| GET    | `/api/admin/stylists/:stylistId` | Get stylist details       | 🔄 TODO        |

### Schedule Management
| Method | Endpoint                                                  | Description              | Status        |
| ------ | --------------------------------------------------------- | ------------------------ | ------------- |
| GET    | `/api/admin/schedules`                                    | List all schedules       | 🔄 TODO        |
| POST   | `/api/admin/schedules/bulk`                               | Create schedules (bulk)  | ✅ Implemented |
| DELETE | `/api/admin/schedules/bulk`                               | Delete schedules (bulk)  | ✅ Implemented |
| GET    | `/api/admin/schedules/:scheduleId`                        | Get schedule details     | 🔄 TODO        |
| PATCH  | `/api/admin/schedules/:scheduleId`                        | Update schedule          | 🔄 TODO        |
| DELETE | `/api/admin/schedules/:scheduleId`                        | Delete schedule          | 🔄 TODO        |
| GET    | `/api/admin/schedules/:scheduleId/time-slots`             | List schedule time slots | 🔄 TODO        |
| POST   | `/api/admin/schedules/:scheduleId/time-slots`             | Create time slot         | ✅ Implemented |
| GET    | `/api/admin/schedules/:scheduleId/time-slots/:timeSlotId` | Get time slot details    | 🔄 TODO        |
| PATCH  | `/api/admin/schedules/:scheduleId/time-slots/:timeSlotId` | Update time slot         | ✅ Implemented |
| DELETE | `/api/admin/schedules/:scheduleId/time-slots/:timeSlotId` | Delete time slot         | ✅ Implemented |

### Time Slot Template Management
| Method | Endpoint                                                   | Description          | Status        |
| ------ | ---------------------------------------------------------- | -------------------- | ------------- |
| GET    | `/api/admin/time-slot-templates`                           | List all templates   | 🔄 TODO        |
| POST   | `/api/admin/time-slot-templates`                           | Create template      | ✅ Implemented |
| GET    | `/api/admin/time-slot-templates/:templateId`               | Get template details | 🔄 TODO        |
| PATCH  | `/api/admin/time-slot-templates/:templateId`               | Update template      | ✅ Implemented |
| DELETE | `/api/admin/time-slot-templates/:templateId`               | Delete template      | ✅ Implemented |
| GET    | `/api/admin/time-slot-templates/:templateId/items`         | List template items  | 🔄 TODO        |
| POST   | `/api/admin/time-slot-templates/:templateId/items`         | Create template item | ✅ Implemented |
| PATCH  | `/api/admin/time-slot-templates/:templateId/items/:itemId` | Update template item | ✅ Implemented |
| DELETE | `/api/admin/time-slot-templates/:templateId/items/:itemId` | Delete template item | ✅ Implemented |

### Booking Management (Admin view)
| Method | Endpoint                         | Description         | Status |
| ------ | -------------------------------- | ------------------- | ------ |
| GET    | `/api/admin/bookings`            | List all bookings   | 🔄 TODO |
| GET    | `/api/admin/bookings/:bookingId` | Get booking details | 🔄 TODO |
| PATCH  | `/api/admin/bookings/:bookingId` | Update booking      | 🔄 TODO |
| DELETE | `/api/admin/bookings/:bookingId` | Cancel booking      | 🔄 TODO |

### Customer Management (Admin view)
| Method | Endpoint                           | Description          | Status |
| ------ | ---------------------------------- | -------------------- | ------ |
| GET    | `/api/admin/customers`             | List all customers   | 🔄 TODO |
| GET    | `/api/admin/customers/:customerId` | Get customer details | 🔄 TODO |
| PATCH  | `/api/admin/customers/:customerId` | Update customer      | 🔄 TODO |

## System Routes
| Method | Endpoint  | Description  | Status        |
| ------ | --------- | ------------ | ------------- |
| GET    | `/health` | Health check | ✅ Implemented |