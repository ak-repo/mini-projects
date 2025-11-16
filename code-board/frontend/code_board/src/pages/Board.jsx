import { Plus, MoreVertical, Filter, Settings } from 'lucide-react';
import Button from '../components/Button';

const columns = [
  {
    id: 'todo',
    title: 'Todo',
    color: 'border-gray-600',
    tasks: [
      {
        id: 1,
        title: 'Design login page UI',
        tags: ['Design', 'UI'],
        assignees: ['AB', 'CD'],
        priority: 'high',
      },
      {
        id: 2,
        title: 'Setup database schema',
        tags: ['Backend', 'Database'],
        assignees: ['EF'],
        priority: 'medium',
      },
      {
        id: 3,
        title: 'Write API documentation',
        tags: ['Documentation'],
        assignees: ['GH'],
        priority: 'low',
      },
    ],
  },
  {
    id: 'in-progress',
    title: 'In Progress',
    color: 'border-[#38bdf8]',
    tasks: [
      {
        id: 4,
        title: 'Implement user authentication',
        tags: ['Frontend', 'Backend'],
        assignees: ['AB', 'EF'],
        priority: 'high',
      },
      {
        id: 5,
        title: 'Create dashboard components',
        tags: ['Frontend', 'React'],
        assignees: ['CD'],
        priority: 'high',
      },
    ],
  },
  {
    id: 'done',
    title: 'Done',
    color: 'border-[#10b981]',
    tasks: [
      {
        id: 6,
        title: 'Setup project repository',
        tags: ['DevOps'],
        assignees: ['AB'],
        priority: 'high',
      },
      {
        id: 7,
        title: 'Configure ESLint and Prettier',
        tags: ['DevOps', 'Config'],
        assignees: ['EF'],
        priority: 'low',
      },
      {
        id: 8,
        title: 'Install dependencies',
        tags: ['Setup'],
        assignees: ['AB'],
        priority: 'medium',
      },
    ],
  },
];

const priorityColors = {
  high: 'bg-red-500/20 text-red-400 border-red-500/30',
  medium: 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30',
  low: 'bg-blue-500/20 text-blue-400 border-blue-500/30',
};

export default function Board() {
  return (
    <div>
      <div className="sticky top-16 -mt-8 -mx-8 px-8 py-6 bg-[#0f172a] border-b border-gray-800 mb-8 z-10">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-100 mb-2">Project Board</h1>
            <p className="text-gray-400">Manage tasks and track progress</p>
          </div>
          <div className="flex items-center gap-3">
            <Button variant="secondary" className="flex items-center gap-2">
              <Filter className="w-4 h-4" />
              Filter
            </Button>
            <Button variant="secondary" className="flex items-center gap-2">
              <Settings className="w-4 h-4" />
              Settings
            </Button>
            <Button variant="primary" className="flex items-center gap-2">
              <Plus className="w-4 h-4" />
              Add Task
            </Button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {columns.map((column) => (
          <div key={column.id} className="flex flex-col">
            <div className={`border-t-4 ${column.color} bg-[#1e293b] rounded-t-2xl px-4 py-3 flex items-center justify-between`}>
              <div className="flex items-center gap-2">
                <h3 className="font-bold text-gray-100">{column.title}</h3>
                <span className="px-2 py-0.5 bg-[#0f172a] rounded-full text-xs text-gray-400">
                  {column.tasks.length}
                </span>
              </div>
              <button className="p-1 hover:bg-[#0f172a] rounded-lg transition-colors">
                <MoreVertical className="w-4 h-4 text-gray-400" />
              </button>
            </div>

            <div className="flex-1 bg-[#1e293b] rounded-b-2xl px-4 py-3 space-y-3">
              {column.tasks.map((task) => (
                <div
                  key={task.id}
                  className="bg-[#0f172a] border border-gray-800 rounded-xl p-4 hover:border-[#38bdf8]/50 transition-colors cursor-pointer group"
                >
                  <div className="flex items-start justify-between mb-3">
                    <h4 className="font-medium text-gray-100 text-sm group-hover:text-[#38bdf8] transition-colors">
                      {task.title}
                    </h4>
                    <button className="opacity-0 group-hover:opacity-100 p-1 hover:bg-[#1e293b] rounded transition-all">
                      <MoreVertical className="w-4 h-4 text-gray-400" />
                    </button>
                  </div>

                  <div className="flex flex-wrap gap-2 mb-3">
                    {task.tags.map((tag) => (
                      <span
                        key={tag}
                        className="px-2 py-1 bg-[#1e293b] border border-gray-700 rounded-lg text-xs text-gray-300"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="flex -space-x-2">
                      {task.assignees.map((assignee, idx) => (
                        <div
                          key={idx}
                          className="w-7 h-7 bg-gradient-to-br from-[#38bdf8] to-[#a855f7] rounded-full flex items-center justify-center text-white text-xs font-semibold border-2 border-[#0f172a]"
                        >
                          {assignee}
                        </div>
                      ))}
                    </div>
                    <span
                      className={`px-2 py-1 rounded-lg text-xs border ${priorityColors[task.priority]}`}
                    >
                      {task.priority}
                    </span>
                  </div>
                </div>
              ))}

              <button className="w-full py-3 border-2 border-dashed border-gray-700 rounded-xl text-gray-400 hover:border-[#38bdf8] hover:text-[#38bdf8] transition-colors flex items-center justify-center gap-2">
                <Plus className="w-4 h-4" />
                Add Task
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
