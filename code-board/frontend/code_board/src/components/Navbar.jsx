import { Search, Bell, Code2 } from 'lucide-react';

export default function Navbar() {
  return (
    <nav className="fixed top-0 left-0 right-0 h-16 bg-[#1e293b] border-b border-gray-800 z-50">
      <div className="h-full px-6 flex items-center justify-between">
        <div className="flex items-center gap-8">
          <div className="flex items-center gap-2">
            <div className="w-10 h-10 bg-gradient-to-br from-[#38bdf8] to-[#a855f7] rounded-xl flex items-center justify-center">
              <Code2 className="w-6 h-6 text-white" />
            </div>
            <span className="text-xl font-bold bg-gradient-to-r from-[#38bdf8] to-[#a855f7] bg-clip-text text-transparent">
              CodeBoard
            </span>
          </div>

          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search projects, boards, repos..."
              className="w-96 pl-10 pr-4 py-2 bg-[#0f172a] border border-gray-700 rounded-xl text-sm focus:outline-none focus:border-[#38bdf8] transition-colors"
            />
          </div>
        </div>

        <div className="flex items-center gap-4">
          <button className="relative p-2 hover:bg-[#0f172a] rounded-lg transition-colors">
            <Bell className="w-5 h-5 text-gray-400" />
            <span className="absolute top-1 right-1 w-2 h-2 bg-[#a855f7] rounded-full"></span>
          </button>

          <div className="flex items-center gap-3 pl-4 border-l border-gray-700">
            <div className="text-right">
              <div className="text-sm font-medium">Alex Developer</div>
              <div className="text-xs text-gray-400">@alexdev</div>
            </div>
            <div className="w-10 h-10 bg-gradient-to-br from-[#38bdf8] to-[#a855f7] rounded-xl flex items-center justify-center text-white font-semibold">
              AD
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
}
