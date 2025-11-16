import { useState } from "react";
import { useAuth } from "../context/context";
import { useNavigate } from "react-router-dom";

export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true);
  const navigate = useNavigate();
  const [form, setForm] = useState({ username: "", password: "", email: "" });
  const { login, register, isAuthenticated } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (isLogin) {
      await login(form.email, form.password);
    } else {
      await register(form);
    }
    if (isAuthenticated) {
      navigate("/");
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex flex-col justify-center px-6 py-12 lg:px-8 relative overflow-hidden">
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-10">
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_1px_1px,rgba(120,119,198,0.15)_1px,transparent_0)] bg-[length:20px_20px]"></div>
      </div>

      {/* Animated Glow Effects */}
      <div className="absolute top-1/4 left-1/4 w-64 h-64 bg-blue-500 rounded-full filter blur-3xl opacity-10 animate-pulse"></div>
      <div className="absolute bottom-1/4 right-1/4 w-64 h-64 bg-purple-500 rounded-full filter blur-3xl opacity-10 animate-pulse delay-1000"></div>

      <div className="sm:mx-auto sm:w-full sm:max-w-md relative z-10">
        <div className="text-center">
          {/* Logo */}
          <div className="flex justify-center mb-6">
            <div className="relative">
              <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/30">
                <span className="text-white font-bold text-xl tracking-wider">
                  CB
                </span>
              </div>
              <div className="absolute inset-0 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl blur-sm opacity-50"></div>
            </div>
          </div>

          <h2 className="text-4xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent tracking-tight">
            CodeBoard
          </h2>
          <p className="mt-3 text-sm text-gray-400 font-light">
            {isLogin
              ? "Welcome back, developer"
              : "Join the developer community"}
          </p>
        </div>

        {/* Toggle Switch */}
        <div className="mt-8 flex justify-center">
          <div className="bg-gray-800 rounded-xl p-1.5 shadow-inner border border-gray-700">
            <button
              onClick={() => setIsLogin(true)}
              className={`px-8 py-3 rounded-lg text-sm font-medium transition-all duration-300 ${
                isLogin
                  ? "bg-gradient-to-r from-blue-600 to-purple-600 text-white shadow-lg shadow-blue-500/25"
                  : "text-gray-400 hover:text-gray-300 hover:bg-gray-700"
              }`}
            >
              Sign In
            </button>
            <button
              onClick={() => setIsLogin(false)}
              className={`px-8 py-3 rounded-lg text-sm font-medium transition-all duration-300 ${
                !isLogin
                  ? "bg-gradient-to-r from-blue-600 to-purple-600 text-white shadow-lg shadow-blue-500/25"
                  : "text-gray-400 hover:text-gray-300 hover:bg-gray-700"
              }`}
            >
              Sign Up
            </button>
          </div>
        </div>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md relative z-10">
        <div className="bg-gray-800/60 backdrop-blur-sm py-10 px-8 rounded-2xl shadow-2xl border border-gray-700/50">
          <form className="space-y-6">
            {!isLogin && (
              <div>
                <label
                  htmlFor="username"
                  className="block text-sm font-medium text-gray-300 mb-2"
                >
                  Full Name
                </label>
                <div className="mt-1">
                  <input
                    id="username"
                    name="username"
                    type="text"
                    value={form.username}
                    onChange={(e) => handleChange(e)}
                    required
                    className="block w-full rounded-xl bg-gray-700/50 border border-gray-600 px-4 py-3.5 text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-300"
                    placeholder="Enter your full name"
                  />
                </div>
              </div>
            )}

            <div>
              <label
                htmlFor="email"
                className="block text-sm font-medium text-gray-300 mb-2"
              >
                Email address
              </label>
              <div className="mt-1">
                <input
                  id="email"
                  name="email"
                  type="email"
                  value={form.email}
                  onChange={(e) => handleChange(e)}
                  required
                  autoComplete="email"
                  className="block w-full rounded-xl bg-gray-700/50 border border-gray-600 px-4 py-3.5 text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-300"
                  placeholder="developer@example.com"
                />
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between mb-2">
                <label
                  htmlFor="password"
                  className="block text-sm font-medium text-gray-300"
                >
                  Password
                </label>
                {isLogin && (
                  <div className="text-sm">
                    <a
                      href="#"
                      className="font-medium text-blue-400 hover:text-blue-300 transition-colors"
                    >
                      Forgot password?
                    </a>
                  </div>
                )}
              </div>
              <div className="mt-1">
                <input
                  onChange={(e) => handleChange(e)}
                  id="password"
                  name="password"
                  type="password"
                  value={form.password}
                  required
                  autoComplete={isLogin ? "current-password" : "new-password"}
                  className="block w-full rounded-xl bg-gray-700/50 border border-gray-600 px-4 py-3.5 text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-300"
                  placeholder={
                    isLogin ? "Enter your password" : "Create a secure password"
                  }
                />
              </div>
            </div>

            <div className="pt-4">
              <button
                onClick={(e) => handleSubmit(e)}
                type="submit"
                className="flex w-full justify-center rounded-xl bg-gradient-to-r from-blue-600 to-purple-600 py-4 px-4 text-sm font-semibold text-white hover:from-blue-500 hover:to-purple-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-900 transition-all duration-300 shadow-lg shadow-blue-500/25 hover:shadow-blue-500/40"
              >
                {isLogin ? "Sign in to CodeBoard" : "Create account"}
              </button>
            </div>
          </form>

          <div className="mt-8 text-center border-t border-gray-700/50 pt-6">
            <p className="text-sm text-gray-400">
              {isLogin ? "Don't have an account?" : "Already have an account?"}{" "}
              <button
                onClick={() => setIsLogin(!isLogin)}
                className="font-medium text-blue-400 hover:text-blue-300 transition-colors"
              >
                {isLogin ? "Join CodeBoard" : "Sign in"}
              </button>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
