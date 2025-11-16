# FreshBox API Services Documentation

This document outlines all API services and their corresponding backend routes for the FreshBox ecommerce application.

## Backend Route Mapping

### Public Routes (No Authentication Required)

#### Authentication

- `POST /api/v1/customer/login` → `customerLogin(loginData)`
- `POST /api/v1/customer/register` → `customerRegister(registerData)`

#### Products

- `GET /api/v1/customer/products` → `getAllProducts()`
- `GET /api/v1/customer/products/:id` → `getProductById(id)`
- `POST /api/v1/customer/products/filter` → `filterProducts(filterData)`

#### Pages

- `GET /api/v1/customer/banners` → `getBanners()`

### Protected Routes (Authentication Required)

All protected routes are prefixed with `/api/v1/customer/auth` and require authentication middleware.

#### Authentication Management

- `POST /api/v1/customer/auth/password-change` → `changePassword(passwordData)`
- `POST /api/v1/customer/auth/send-otp` → `sendOTP(otpData)`
- `POST /api/v1/customer/auth/verify-otp` → `verifyOTP(otpData)`
- `GET /api/v1/customer/auth/me` → `getMe()`

#### Profile Management

- `GET /api/v1/customer/auth/profile` → `getProfile()`
- `GET /api/v1/customer/auth/profile/address` → `getAddress()`
- `PATCH /api/v1/customer/auth/profile/address` → `updateAddress(addressData)`

#### Cart Management

- `GET /api/v1/customer/auth/cart` → `getUserCart()`
- `POST /api/v1/customer/auth/cart` → `addToCart(cartData)`
- `PATCH /api/v1/customer/auth/cart` → `updateCartQuantity(updateData)`
- `DELETE /api/v1/customer/auth/cart` → `removeFromCart(removeData)`

#### Checkout Process

- `GET /api/v1/customer/auth/checkout` → `getCheckoutSummary()`
- `POST /api/v1/customer/auth/checkout` → `processCheckout(checkoutData)`

#### Order Management

- `GET /api/v1/customer/auth/orders` → `getOrdersByCustomerId()`
- `GET /api/v1/customer/auth/orders/:id` → `getOrderById(id)`
- `POST /api/v1/customer/auth/orders/cancel` → `cancelOrder(cancelData)`
- `GET /api/v1/customer/auth/orders/cancel-response/:id` → `getCancellationResponse(id)`

#### Wishlist Management

- `GET /api/v1/customer/auth/wishlist` → `getWishlist()`
- `POST /api/v1/customer/auth/wishlist/:id` → `addToWishlist(productId)`
- `DELETE /api/v1/customer/auth/wishlist/:id` → `removeFromWishlist(productId)`

#### Reviews

- `POST /api/v1/customer/auth/review` → `addReview(reviewData)`

## Service Files Structure

```
src/api/
├── api.js                    # Base API configuration
├── index.js                  # Main export file
├── services/
│   ├── index.js             # Services export file
│   ├── authService.js        # Authentication services
│   ├── productService.js     # Product-related services
│   ├── cartService.js        # Cart management services
│   ├── checkoutService.js    # Checkout process services
│   ├── orderService.js       # Order management services
│   ├── wishlistService.js    # Wishlist services
│   ├── profileService.js     # User profile services
│   ├── reviewService.js      # Review services
│   ├── pagesService.js       # Pages/banners services
│   └── contactService.js     # Contact form services
└── API_DOCUMENTATION.md      # This documentation
```

## Usage Examples

### Authentication

```javascript
import { customerLogin, customerRegister } from "../api/services/authService";

// Login
const loginData = { email: "user@example.com", password: "password" };
const user = await customerLogin(loginData);

// Register
const registerData = {
  name: "John",
  email: "john@example.com",
  password: "password",
};
const newUser = await customerRegister(registerData);
```

### Products

```javascript
import {
  getAllProducts,
  getProductById,
  filterProducts,
} from "../api/services/productService";

// Get all products
const products = await getAllProducts();

// Get specific product
const product = await getProductById(1);

// Filter products
const filteredProducts = await filterProducts({
  category: "smoothies",
  price_range: "0-100",
});
```

### Cart Management

```javascript
import {
  getUserCart,
  addToCart,
  updateCartQuantity,
  removeFromCart,
} from "../api/services/cartService";

// Get user's cart
const cart = await getUserCart();

// Add item to cart
await addToCart({ product_id: 1, quantity: 2 });

// Update quantity
await updateCartQuantity({ cart_item_id: 1, quantity: 3 });

// Remove item
await removeFromCart({ cart_item_id: 1 });
```

### Checkout Process

```javascript
import {
  getCheckoutSummary,
  processCheckout,
} from "../api/services/checkoutService";

// Get checkout summary
const summary = await getCheckoutSummary();

// Process checkout
const order = await processCheckout({
  shipping_address: {
    /* address data */
  },
  payment_method: "card",
  payment_details: {
    /* payment data */
  },
});
```

## Error Handling

All services include proper error handling:

```javascript
try {
  const data = await getAllProducts();
  setProducts(data);
} catch (error) {
  console.error("Error fetching products:", error);
  // Handle error appropriately
  setError("Failed to load products");
}
```

## Authentication

The API automatically handles authentication by:

1. Reading the token from localStorage
2. Adding it to the Authorization header
3. Including it in all protected route requests

## Response Format

All API responses follow this format:

```javascript
{
  success: true,
  data: { /* actual data */ },
  message: "Success message"
}
```

The response interceptor automatically extracts the `data` field, so you receive the actual data directly.

## Development Notes

- All services use async/await for promise handling
- Services include fallback data for development/demo purposes
- Error messages are logged to console for debugging
- Services are designed to be easily testable and mockable
- All routes match exactly with the backend Gin routes provided

