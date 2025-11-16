import { Link } from 'react-router-dom';
import { Home, Code2 } from 'lucide-react';
import Button from '../components/Button';

export default function NotFound() {
  return (
    <div className="min-h-screen bg-[#0f172a] relative overflow-hidden flex items-center justify-center">
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute w-96 h-96 bg-[#38bdf8]/10 rounded-full blur-3xl top-1/4 left-1/4 animate-pulse"></div>
        <div className="absolute w-96 h-96 bg-[#a855f7]/10 rounded-full blur-3xl bottom-1/4 right-1/4 animate-pulse delay-700"></div>

        {[...Array(30)].map((_, i) => (
          <div
            key={i}
            className="absolute text-gray-800 font-mono text-lg opacity-20"
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
              animation: `float ${3 + Math.random() * 6}s infinite ease-in-out`,
              animationDelay: `${Math.random() * 3}s`,
            }}
          >
            {['{', '}', '<', '>', '/', '(', ')', ';', '=', '[', ']'][Math.floor(Math.random() * 11)]}
          </div>
        ))}
      </div>

      <div className="relative z-10 text-center px-6">
        <div className="mb-8">
          <div className="inline-flex items-center gap-3 mb-6">
            <div className="w-16 h-16 bg-gradient-to-br from-[#38bdf8] to-[#a855f7] rounded-2xl flex items-center justify-center animate-pulse">
              <Code2 className="w-10 h-10 text-white" />
            </div>
          </div>

          <div className="relative mb-6">
            <h1 className="text-[150px] md:text-[200px] font-bold bg-gradient-to-r from-[#38bdf8] via-[#a855f7] to-[#38bdf8] bg-clip-text text-transparent leading-none animate-gradient">
              404
            </h1>
            <div className="absolute inset-0 blur-3xl bg-gradient-to-r from-[#38bdf8]/20 via-[#a855f7]/20 to-[#38bdf8]/20"></div>
          </div>

          <h2 className="text-3xl md:text-4xl font-bold text-gray-100 mb-4">Page Not Found</h2>
          <p className="text-gray-400 text-lg mb-8 max-w-md mx-auto">
            The page you're looking for doesn't exist or has been moved to another location.
          </p>

          <div className="flex items-center justify-center gap-4">
            <Link to="/dashboard">
              <Button variant="primary" className="flex items-center gap-2">
                <Home className="w-5 h-5" />
                Back to Dashboard
              </Button>
            </Link>
            <Link to="/login">
              <Button variant="secondary">
                Go to Login
              </Button>
            </Link>
          </div>
        </div>

        <div className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-4 max-w-2xl mx-auto">
          {[
            { label: 'Dashboard', path: '/dashboard' },
            { label: 'Boards', path: '/boards/main' },
            { label: 'Repositories', path: '/repos/demo-project' },
          ].map((link) => (
            <Link key={link.path} to={link.path}>
              <div className="bg-[#1e293b] border border-gray-800 rounded-xl p-4 hover:border-[#38bdf8] hover:scale-105 transition-all cursor-pointer">
                <span className="text-gray-300 hover:text-[#38bdf8] font-medium">{link.label}</span>
              </div>
            </Link>
          ))}
        </div>
      </div>

      <style>{`
        @keyframes float {
          0%, 100% {
            transform: translateY(0) translateX(0) rotate(0deg);
            opacity: 0.2;
          }
          50% {
            transform: translateY(-30px) translateX(20px) rotate(180deg);
            opacity: 0.3;
          }
        }
        @keyframes gradient {
          0%, 100% { background-position: 0% 50%; }
          50% { background-position: 100% 50%; }
        }
        .animate-gradient {
          background-size: 200% auto;
          animation: gradient 3s ease infinite;
        }
      `}</style>
    </div>
  );
}
