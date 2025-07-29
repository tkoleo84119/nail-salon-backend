# API Routes Documentation

This document lists all API routes in the nail salon backend system, organized by public/customer routes and admin routes.

## Public/Customer Routes

### Authentication
| Method | Endpoint                  | Description                | Status        |
| ------ | ------------------------- | -------------------------- | ------------- |
| POST   | `/api/auth/line/login`    | Customer LINE login        | ✅ Implemented |
| POST   | `/api/auth/line/register` | Customer LINE registration | ✅ Implemented |
| POST   | `/api/auth/token/refresh` | Refresh access token       | ✅ Implemented |

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
| GET    | `/api/bookings/:bookingId`        | Get booking details | ✅ Implemented |
| PATCH  | `/api/bookings/:bookingId`        | Update my booking   | ✅ Implemented |
| PATCH  | `/api/bookings/:bookingId/cancel` | Cancel my booking   | ✅ Implemented |

### Browse Stores (Read-only)
| Method | Endpoint                        | Description         | Status        |
| ------ | ------------------------------- | ------------------- | ------------- |
| GET    | `/api/stores`                   | List stores         | ✅ Implemented |
| GET    | `/api/stores/:storeId`          | Get store details   | ✅ Implemented |
| GET    | `/api/stores/:storeId/stylists` | List store stylists | ✅ Implemented |
| GET    | `/api/stores/:storeId/services` | List store services | ✅ Implemented |

### Browse Schedules & Time Slots (Read-only)
| Method | Endpoint                                             | Description                  | Status        |
| ------ | ---------------------------------------------------- | ---------------------------- | ------------- |
| GET    | `/api/stores/:storeId/stylists/:stylistId/schedules` | List store stylist schedules | ✅ Implemented |
| GET    | `/api/schedules/:scheduleId/time-slots`              | List available time slots    | ✅ Implemented |

## Admin Routes

### Authentication
| Method | Endpoint                        | Description          | Status        |
| ------ | ------------------------------- | -------------------- | ------------- |
| POST   | `/api/admin/auth/login`         | Staff login          | ✅ Implemented |
| POST   | `/api/admin/auth/token/refresh` | Refresh access token | ✅ Implemented |

### Staff Management
| Method | Endpoint                    | Description       | Status        |
| ------ | --------------------------- | ----------------- | ------------- |
| GET    | `/api/admin/staff`          | List all staff    | ✅ Implemented |
| POST   | `/api/admin/staff`          | Create staff      | ✅ Implemented |
| GET    | `/api/admin/staff/me`       | Get my profile    | ✅ Implemented |
| PATCH  | `/api/admin/staff/me`       | Update my profile | ✅ Implemented |
| GET    | `/api/admin/staff/:staffId` | Get staff details | ✅ Implemented |
| PATCH  | `/api/admin/staff/:staffId` | Update staff      | ✅ Implemented |

### Staff Store Access
| Method | Endpoint                                      | Description                | Status        |
| ------ | --------------------------------------------- | -------------------------- | ------------- |
| GET    | `/api/admin/staff/:staffId/store-access`      | List staff store access    | ✅ Implemented |
| POST   | `/api/admin/staff/:staffId/store-access`      | Grant store access         | ✅ Implemented |
| DELETE | `/api/admin/staff/:staffId/store-access/bulk` | Revoke store access (bulk) | ✅ Implemented |

### Store Management
| Method | Endpoint                     | Description       | Status        |
| ------ | ---------------------------- | ----------------- | ------------- |
| GET    | `/api/admin/stores`          | List all stores   | ✅ Implemented |
| POST   | `/api/admin/stores`          | Create store      | ✅ Implemented |
| GET    | `/api/admin/stores/:storeId` | Get store details | ✅ Implemented |
| PATCH  | `/api/admin/stores/:storeId` | Update store      | ✅ Implemented |

### Service Management
| Method | Endpoint                                         | Description                | Status        |
| ------ | ------------------------------------------------ | -------------------------- | ------------- |
| GET    | `/api/admin/stores/:storeId/services`            | List all services in store | ✅ Implemented |
| POST   | `/api/admin/stores/:storeId/services`            | Create service             | ✅ Implemented |
| GET    | `/api/admin/stores/:storeId/services/:serviceId` | Get service details        | ✅ Implemented |
| PATCH  | `/api/admin/stores/:storeId/services/:serviceId` | Update service             | ✅ Implemented |

### Stylist Management
| Method | Endpoint                              | Description               | Status        |
| ------ | ------------------------------------- | ------------------------- | ------------- |
| GET    | `/api/admin/stores/:storeId/stylists` | List all stylists         | ✅ Implemented |
| POST   | `/api/admin/stylists/me`              | Create my stylist profile | ✅ Implemented |
| PATCH  | `/api/admin/stylists/me`              | Update my stylist profile | ✅ Implemented |

### Schedule Management
| Method | Endpoint                                                  | Description             | Status        |
| ------ | --------------------------------------------------------- | ----------------------- | ------------- |
| GET    | `/api/admin/stores/:storeId/schedules`                    | List all schedules      | ✅ Implemented |
| POST   | `/api/admin/stores/:storeId/schedules/bulk`               | Create schedules (bulk) | ✅ Implemented |
| DELETE | `/api/admin/stores/:storeId/schedules/bulk`               | Delete schedules (bulk) | ✅ Implemented |
| GET    | `/api/admin/stores/:storeId/schedules/:scheduleId`        | Get schedule details    | ✅ Implemented |
| POST   | `/api/admin/schedules/:scheduleId/time-slots`             | Create time slot        | ✅ Implemented |
| PATCH  | `/api/admin/schedules/:scheduleId/time-slots/:timeSlotId` | Update time slot        | ✅ Implemented |
| DELETE | `/api/admin/schedules/:scheduleId/time-slots/:timeSlotId` | Delete time slot        | ✅ Implemented |

### Time Slot Template Management
| Method | Endpoint                                                   | Description          | Status        |
| ------ | ---------------------------------------------------------- | -------------------- | ------------- |
| GET    | `/api/admin/time-slot-templates`                           | List all templates   | ✅ Implemented |
| POST   | `/api/admin/time-slot-templates`                           | Create template      | ✅ Implemented |
| GET    | `/api/admin/time-slot-templates/:templateId`               | Get template details | ✅ Implemented |
| PATCH  | `/api/admin/time-slot-templates/:templateId`               | Update template      | ✅ Implemented |
| DELETE | `/api/admin/time-slot-templates/:templateId`               | Delete template      | ✅ Implemented |
| POST   | `/api/admin/time-slot-templates/:templateId/items`         | Create template item | ✅ Implemented |
| PATCH  | `/api/admin/time-slot-templates/:templateId/items/:itemId` | Update template item | ✅ Implemented |
| DELETE | `/api/admin/time-slot-templates/:templateId/items/:itemId` | Delete template item | ✅ Implemented |

### Booking Management (Admin view)
| Method | Endpoint                                                | Description       | Status        |
| ------ | ------------------------------------------------------- | ----------------- | ------------- |
| POST   | `/api/admin/stores/:storeId/bookings`                   | Create booking    | ✅ Implemented |
| GET    | `/api/admin/stores/:storeId/bookings`                   | List all bookings | ✅ Implemented |
| PATCH  | `/api/admin/stores/:storeId/bookings/:bookingId`        | Update booking    | ✅ Implemented |
| PATCH  | `/api/admin/stores/:storeId/bookings/:bookingId/cancel` | Cancel booking    | ✅ Implemented |

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