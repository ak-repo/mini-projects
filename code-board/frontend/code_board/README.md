customer-frontend/
│
├── public/ # Static assets (favicon, logos)
│
├── src/
│ ├── assets/ # Images, icons, and static resources
│ │ └── logo.svg
│ │
│ ├── components/ # Reusable UI building blocks
│ │ ├── Navbar.jsx
│ │ ├── Footer.jsx
│ │ ├── ProductCard.jsx
│ │ ├── Loader.jsx
│ │ └── RatingStars.jsx
│ │
│ ├── pages/ # Route-based views
│ │ ├── Home.jsx
│ │ ├── Products.jsx
│ │ ├── ProductDetail.jsx
│ │ ├── Cart.jsx
│ │ ├── Checkout.jsx
│ │ ├── Orders.jsx
│ │ ├── Login.jsx
│ │ ├── Signup.jsx
│ │ └── Profile.jsx
│ │
│ ├── layouts/ # Common page layouts
│ │ ├── MainLayout.jsx # Navbar + Footer wrapper
│ │ └── AuthLayout.jsx # For login/signup pages
│ │
│ ├── context/ # Context API (auth/cart etc.)
│ │ ├── AuthContext.jsx
│ │ └── CartContext.jsx
│ │
│ ├── hooks/ # Custom React hooks
│ │ ├── useAuth.js
│ │ └── useFetch.js
│ │
│ ├── services/ # API layer for backend communication
│ │ ├── api.js # axios/fetch setup
│ │ ├── productService.js
│ │ ├── authService.js
│ │ └── orderService.js
│ │
│ ├── utils/ # Helpers (formatting, validation, etc.)
│ │ ├── formatCurrency.js
│ │ ├── handleError.js
│ │ └── constants.js
│ │
│ ├── router/ # All routes managed here
│ │ └── AppRouter.jsx
│ │
│ ├── App.jsx # Root component
│ ├── main.jsx # Entry point
│ └── index.css # Tailwind / global styles
│
├── .env # API base URL etc. (REACT_APP_API_URL)
├── package.json
├── tailwind.config.js
└── vite.config.js
