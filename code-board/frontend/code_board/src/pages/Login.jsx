import { Link } from 'react-router-dom';
import { Code2, Mail, Lock } from 'lucide-react';
import Button from '../components/Button';

export default function Login() {
  return (
    <div className="min-h-screen bg-[#0f172a] relative overflow-hidden flex items-center justify-center">
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute w-96 h-96 bg-[#38bdf8]/10 rounded-full blur-3xl -top-20 -left-20 animate-pulse"></div>
        <div className="absolute w-96 h-96 bg-[#a855f7]/10 rounded-full blur-3xl -bottom-20 -right-20 animate-pulse delay-700"></div>
        {[...Array(20)].map((_, i) => (
          <div
            key={i}
            className="absolute text-gray-800 font-mono text-xs opacity-20"
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
              animation: `float ${5 + Math.random() * 10}s infinite`,
            }}
          >
            {['const', 'function', 'return', 'import', 'export'][Math.floor(Math.random() * 5)]}
          </div>
        ))}
      </div>

      <div className="relative z-10 w-full max-w-md px-6">
        <div className="text-center mb-8">
          <div className="inline-flex items-center gap-3 mb-4">
            <div className="w-14 h-14 bg-gradient-to-br from-[#38bdf8] to-[#a855f7] rounded-2xl flex items-center justify-center">
              <Code2 className="w-8 h-8 text-white" />
            </div>
            <span className="text-3xl font-bold bg-gradient-to-r from-[#38bdf8] to-[#a855f7] bg-clip-text text-transparent">
              CodeBoard
            </span>
          </div>
          <h1 className="text-2xl font-bold text-gray-100 mb-2">Welcome Back</h1>
          <p className="text-gray-400">Sign in to continue to your projects</p>
        </div>

        <div className="bg-[#1e293b] border border-gray-800 rounded-2xl p-8 shadow-2xl backdrop-blur-sm">
          <form className="space-y-5">
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Email</label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                <input
                  type="email"
                  placeholder="you@example.com"
                  className="w-full pl-10 pr-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Password</label>
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                <input
                  type="password"
                  placeholder="••••••••"
                  className="w-full pl-10 pr-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                />
              </div>
            </div>

            <div className="flex items-center justify-between text-sm">
              <label className="flex items-center gap-2 text-gray-400 cursor-pointer">
                <input type="checkbox" className="w-4 h-4 rounded border-gray-700" />
                Remember me
              </label>
              <a href="#" className="text-[#38bdf8] hover:text-[#a855f7] transition-colors">
                Forgot password?
              </a>
            </div>

            <Link to="/dashboard">
              <Button variant="primary" className="w-full">
                Sign In
              </Button>
            </Link>
          </form>

          <div className="mt-6 text-center text-sm text-gray-400">
            Don't have an account?{' '}
            <Link to="/register" className="text-[#38bdf8] hover:text-[#a855f7] transition-colors font-medium">
              Create Account
            </Link>
          </div>
        </div>
      </div>

      <style>{`
        @keyframes float {
          0%, 100% { transform: translateY(0) rotate(0deg); }
          50% { transform: translateY(-20px) rotate(10deg); }
        }
      `}</style>
    </div>
  );
}
