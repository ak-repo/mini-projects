import { Plus, Star, GitBranch, Users, Clock } from 'lucide-react';
import Card from '../components/Card';
import Button from '../components/Button';

const projects = [
  {
    id: 1,
    title: 'E-Commerce Platform',
    description: 'Full-stack online shopping platform with payment integration',
    tags: ['React', 'Node.js', 'PostgreSQL'],
    stars: 24,
    commits: 156,
    members: 3,
    updated: '2 hours ago',
    color: 'from-[#38bdf8] to-[#06b6d4]',
  },
  {
    id: 2,
    title: 'Task Management App',
    description: 'Kanban-style project management tool for teams',
    tags: ['React', 'Firebase', 'TailwindCSS'],
    stars: 18,
    commits: 89,
    members: 2,
    updated: '5 hours ago',
    color: 'from-[#a855f7] to-[#8b5cf6]',
  },
  {
    id: 3,
    title: 'AI Chat Assistant',
    description: 'Intelligent chatbot with natural language processing',
    tags: ['Python', 'OpenAI', 'FastAPI'],
    stars: 42,
    commits: 203,
    members: 5,
    updated: '1 day ago',
    color: 'from-[#f59e0b] to-[#d97706]',
  },
  {
    id: 4,
    title: 'Social Media Dashboard',
    description: 'Analytics platform for social media metrics',
    tags: ['Vue.js', 'Express', 'MongoDB'],
    stars: 31,
    commits: 127,
    members: 4,
    updated: '3 days ago',
    color: 'from-[#10b981] to-[#059669]',
  },
  {
    id: 5,
    title: 'Weather Forecast App',
    description: 'Real-time weather data with beautiful visualizations',
    tags: ['React Native', 'Redux', 'API'],
    stars: 15,
    commits: 64,
    members: 2,
    updated: '1 week ago',
    color: 'from-[#06b6d4] to-[#0891b2]',
  },
  {
    id: 6,
    title: 'Portfolio Generator',
    description: 'Create stunning developer portfolios in minutes',
    tags: ['Next.js', 'TypeScript', 'MDX'],
    stars: 28,
    commits: 98,
    members: 1,
    updated: '2 weeks ago',
    color: 'from-[#ec4899] to-[#db2777]',
  },
];

const stats = [
  { label: 'Total Projects', value: '24', icon: GitBranch, color: 'text-[#38bdf8]' },
  { label: 'Active Boards', value: '12', icon: Users, color: 'text-[#a855f7]' },
  { label: 'Total Stars', value: '186', icon: Star, color: 'text-[#f59e0b]' },
  { label: 'This Month', value: '8', icon: Clock, color: 'text-[#10b981]' },
];

export default function Dashboard() {
  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-100 mb-2">Dashboard</h1>
          <p className="text-gray-400">Welcome back, Alex! Here's what's happening with your projects.</p>
        </div>
        <Button variant="primary" className="flex items-center gap-2">
          <Plus className="w-5 h-5" />
          New Project
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <Card key={stat.label}>
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-400 text-sm mb-1">{stat.label}</p>
                  <p className="text-3xl font-bold text-gray-100">{stat.value}</p>
                </div>
                <Icon className={`w-10 h-10 ${stat.color}`} />
              </div>
            </Card>
          );
        })}
      </div>

      <div className="mb-6 flex items-center gap-4">
        <h2 className="text-xl font-bold text-gray-100">Your Projects</h2>
        <div className="flex gap-2">
          <button className="px-4 py-2 bg-[#1e293b] border border-[#38bdf8] text-[#38bdf8] rounded-lg text-sm transition-colors">
            All
          </button>
          <button className="px-4 py-2 bg-[#1e293b] border border-gray-700 text-gray-400 rounded-lg text-sm hover:border-[#38bdf8] hover:text-[#38bdf8] transition-colors">
            Active
          </button>
          <button className="px-4 py-2 bg-[#1e293b] border border-gray-700 text-gray-400 rounded-lg text-sm hover:border-[#38bdf8] hover:text-[#38bdf8] transition-colors">
            Archived
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
        {projects.map((project) => (
          <Card key={project.id} hover>
            <div className="mb-4">
              <div className={`w-12 h-12 bg-gradient-to-br ${project.color} rounded-xl mb-3`}></div>
              <h3 className="text-lg font-bold text-gray-100 mb-2">{project.title}</h3>
              <p className="text-sm text-gray-400 line-clamp-2">{project.description}</p>
            </div>

            <div className="flex flex-wrap gap-2 mb-4">
              {project.tags.map((tag) => (
                <span
                  key={tag}
                  className="px-3 py-1 bg-[#0f172a] border border-gray-700 rounded-lg text-xs text-gray-300"
                >
                  {tag}
                </span>
              ))}
            </div>

            <div className="flex items-center justify-between text-sm text-gray-400 pt-4 border-t border-gray-800">
              <div className="flex items-center gap-4">
                <div className="flex items-center gap-1">
                  <Star className="w-4 h-4" />
                  <span>{project.stars}</span>
                </div>
                <div className="flex items-center gap-1">
                  <GitBranch className="w-4 h-4" />
                  <span>{project.commits}</span>
                </div>
                <div className="flex items-center gap-1">
                  <Users className="w-4 h-4" />
                  <span>{project.members}</span>
                </div>
              </div>
              <div className="flex items-center gap-1">
                <Clock className="w-4 h-4" />
                <span>{project.updated}</span>
              </div>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
