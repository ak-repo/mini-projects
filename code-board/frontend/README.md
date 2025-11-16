# user-frontend

codeboard-frontend/
│
├─ 📁 public/
│ ├─ index.html
│ ├─ favicon.ico
│ └─ logo.svg
│
├─ 📁 src/
│ ├─ 📁 assets/
│ │ ├─ images/
│ │ ├─ icons/
│ │ └─ illustrations/
│ │
│ ├─ 📁 components/
│ │ ├─ common/
│ │ │ ├─ Navbar.jsx
│ │ │ ├─ Sidebar.jsx
│ │ │ ├─ Footer.jsx
│ │ │ ├─ Button.jsx
│ │ │ └─ Card.jsx
│ │ │
│ │ ├─ layout/
│ │ │ ├─ MainLayout.jsx # For user dashboard layout
│ │ │ ├─ AdminLayout.jsx # For admin dashboard layout
│ │ │ └─ AuthLayout.jsx # For login/register layout
│ │ │
│ │ └─ protected/
│ │ └─ ProtectedRoute.jsx
│ │
│ ├─ 📁 context/
│ │ └─ AuthContext.jsx
│ │
│ ├─ 📁 hooks/
│ │ └─ useAuth.js # custom hook for role-based logic
│ │
│ ├─ 📁 pages/
│ │ ├─ 📁 auth/
│ │ │ ├─ Login.jsx
│ │ │ └─ Register.jsx
│ │ │
│ │ ├─ 📁 user/
│ │ │ ├─ Dashboard.jsx # /user/dashboard
│ │ │ ├─ Boards.jsx # /user/boards
│ │ │ ├─ BoardView.jsx # /user/boards/:id
│ │ │ ├─ Repository.jsx # /user/repos/:id
│ │ │ └─ Settings.jsx # /user/settings
│ │ │
│ │ ├─ 📁 admin/
│ │ │ ├─ Dashboard.jsx # /admin/dashboard
│ │ │ ├─ ManageUsers.jsx # /admin/users
│ │ │ ├─ ManageProjects.jsx # /admin/projects
│ │ │ └─ AdminSettings.jsx # /admin/settings
│ │ │
│ │ └─ NotFound.jsx # 404 Page
│ │
│ ├─ 📁 routes/
│ │ ├─ AppRoutes.jsx # Main route configuration
│ │ └─ AdminRoutes.jsx # Optional separate admin route file
│ │
│ ├─ 📁 styles/
│ │ ├─ index.css # Tailwind base imports
│ │ └─ tailwind.css # Custom CSS utilities
│ │
│ ├─ 📁 utils/
│ │ ├─ constants.js
│ │ ├─ helpers.js
│ │ └─ mockData.js # Fake data for UI placeholders
│ │
│ ├─ App.jsx
│ ├─ index.jsx
│ └─ main.jsx # If using Vite
│
├─ .gitignore
├─ package.json
├─ tailwind.config.js
├─ postcss.config.js
└─ vite.config.js # (if using Vite)
