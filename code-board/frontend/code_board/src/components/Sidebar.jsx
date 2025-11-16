import { Link, useLocation } from 'react-router-dom';
import { LayoutDashboard, Trello, GitBranch, Settings } from 'lucide-react';

const menuItems = [
  { icon: LayoutDashboard, label: 'Dashboard', path: '/dashboard' },
  { icon: Trello, label: 'Boards', path: '/boards/main' },
  { icon: GitBranch, label: 'Repositories', path: '/repos/demo-project' },
  { icon: Settings, label: 'Settings', path: '/settings' },
];

export default function Sidebar() {
  const location = useLocation();

  return (
    <aside className="fixed left-0 top-16 w-64 h-[calc(100vh-4rem)] bg-[#1e293b] border-r border-gray-800">
      <nav className="p-4 space-y-2">
        {menuItems.map((item) => {
          const Icon = item.icon;
          const isActive = location.pathname.startsWith(item.path);

          return (
            <Link
              key={item.path}
              to={item.path}
              className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${
                isActive
                  ? 'bg-gradient-to-r from-[#38bdf8]/20 to-[#a855f7]/20 border border-[#38bdf8]/30 text-[#38bdf8]'
                  : 'text-gray-400 hover:bg-[#0f172a] hover:text-gray-200'
              }`}
            >
              <Icon className="w-5 h-5" />
              <span className="font-medium">{item.label}</span>
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
