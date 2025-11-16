import { Route, Routes } from "react-router-dom";
import AuthPage from "./pages/AuthPage";
import HomePage from "./pages/HomePage";

import Layout from "./components/Layout";

import Dashboard from "./pages/Dashboard";
import Board from "./pages/Board";
import Repository from "./pages/Repository";
import Settings from "./pages/Settings";
import NotFound from "./pages/NotFound";

function AppRoutes() {
  return (
    <>
      <AuthRoutes />
    </>
  );
}

export default AppRoutes;

function AuthRoutes() {
  return (
    <>

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/auth" element={<AuthPage />} />


      

        <Route element={<Layout />}>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/boards/:id" element={<Board />} />
          <Route path="/repos/:id" element={<Repository />} />
          <Route path="/settings" element={<Settings />} />
        </Route>

        <Route path="*" element={<NotFound />} />
      </Routes>
    </>
  );
}
