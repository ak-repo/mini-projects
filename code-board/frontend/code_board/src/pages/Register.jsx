import { Link } from 'react-router-dom';
import { Code2, User, Mail, Lock } from 'lucide-react';
import Button from '../components/Button';

export default function Register() {
  return (
    <div className="min-h-screen bg-[#0f172a] relative overflow-hidden">
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute w-96 h-96 bg-[#38bdf8]/10 rounded-full blur-3xl top-1/4 -left-20 animate-pulse"></div>
        <div className="absolute w-96 h-96 bg-[#a855f7]/10 rounded-full blur-3xl bottom-1/4 -right-20 animate-pulse delay-700"></div>
      </div>

      <div className="relative z-10 min-h-screen flex">
        <div className="flex-1 flex items-center justify-center p-8">
          <div className="w-full max-w-md">
            <div className="text-center mb-8">
              <div className="inline-flex items-center gap-3 mb-4">
                <div className="w-14 h-14 bg-gradient-to-br from-[#38bdf8] to-[#a855f7] rounded-2xl flex items-center justify-center">
                  <Code2 className="w-8 h-8 text-white" />
                </div>
                <span className="text-3xl font-bold bg-gradient-to-r from-[#38bdf8] to-[#a855f7] bg-clip-text text-transparent">
                  CodeBoard
                </span>
              </div>
              <h1 className="text-2xl font-bold text-gray-100 mb-2">Create Account</h1>
              <p className="text-gray-400">Join thousands of developers</p>
            </div>

            <div className="bg-[#1e293b] border border-gray-800 rounded-2xl p-8 shadow-2xl">
              <form className="space-y-5">
                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-2">Username</label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                    <input
                      type="text"
                      placeholder="johndoe"
                      className="w-full pl-10 pr-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                    />
                  </div>
                </div>

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

                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-2">Confirm Password</label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                    <input
                      type="password"
                      placeholder="••••••••"
                      className="w-full pl-10 pr-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                    />
                  </div>
                </div>

                <label className="flex items-start gap-2 text-sm text-gray-400 cursor-pointer">
                  <input type="checkbox" className="w-4 h-4 mt-0.5 rounded border-gray-700" />
                  <span>
                    I agree to the{' '}
                    <a href="#" className="text-[#38bdf8] hover:text-[#a855f7] transition-colors">
                      Terms of Service
                    </a>{' '}
                    and{' '}
                    <a href="#" className="text-[#38bdf8] hover:text-[#a855f7] transition-colors">
                      Privacy Policy
                    </a>
                  </span>
                </label>

                <Link to="/dashboard">
                  <Button variant="primary" className="w-full">
                    Create Account
                  </Button>
                </Link>
              </form>

              <div className="mt-6 text-center text-sm text-gray-400">
                Already have an account?{' '}
                <Link to="/login" className="text-[#38bdf8] hover:text-[#a855f7] transition-colors font-medium">
                  Sign In
                </Link>
              </div>
            </div>
          </div>
        </div>

        <div className="hidden lg:flex flex-1 items-center justify-center p-8">
          <div className="max-w-lg">
            <div className="bg-gradient-to-br from-[#38bdf8]/20 to-[#a855f7]/20 border border-[#38bdf8]/30 rounded-3xl p-12 backdrop-blur-sm">
              <h2 className="text-4xl font-bold text-gray-100 mb-6">
                Start Building
                <br />
                <span className="bg-gradient-to-r from-[#38bdf8] to-[#a855f7] bg-clip-text text-transparent">
                  Amazing Projects
                </span>
              </h2>
              <ul className="space-y-4 text-gray-300">
                <li className="flex items-center gap-3">
                  <div className="w-6 h-6 bg-[#38bdf8] rounded-full flex items-center justify-center text-white text-xs">✓</div>
                  Unlimited repositories and boards
                </li>
                <li className="flex items-center gap-3">
                  <div className="w-6 h-6 bg-[#38bdf8] rounded-full flex items-center justify-center text-white text-xs">✓</div>
                  Collaborate with your team
                </li>
                <li className="flex items-center gap-3">
                  <div className="w-6 h-6 bg-[#38bdf8] rounded-full flex items-center justify-center text-white text-xs">✓</div>
                  Advanced project management tools
                </li>
                <li className="flex items-center gap-3">
                  <div className="w-6 h-6 bg-[#38bdf8] rounded-full flex items-center justify-center text-white text-xs">✓</div>
                  Track progress with kanban boards
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
