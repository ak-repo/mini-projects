import { ChevronRight, Plus, GitBranch, Bug, Users, Activity, Star, Clock } from 'lucide-react';

export default function HomePage() {
  const user = { username: "Developer" };

  const recentBoards = [
    {
      id: 1,
      name: "Frontend Development",
      color: "from-blue-500 to-blue-600",
      updated: "2 hours ago",
    },
    { id: 2, name: "Backend API", color: "from-emerald-500 to-emerald-600", updated: "1 day ago" },
    { id: 3, name: "Mobile App Launch", color: "from-amber-500 to-amber-600", updated: "3 days ago" },
  ];

  const recentRepos = [
    { id: 1, name: "codeboard-web", language: "TypeScript", stars: 12 },
    { id: 2, name: "api-gateway", language: "Go", stars: 8 },
    { id: 3, name: "mobile-app", language: "React Native", stars: 5 },
  ];

  const quickActions = [
    { icon: Plus, label: "New Board", gradient: "from-blue-500 to-cyan-500" },
    { icon: GitBranch, label: "New Repo", gradient: "from-emerald-500 to-teal-500" },
    { icon: Bug, label: "New Issue", gradient: "from-orange-500 to-red-500" },
    { icon: Users, label: "Invite Team", gradient: "from-violet-500 to-purple-500" },
  ];

  return (
    <div className="min-h-screen bg-slate-950 relative overflow-hidden">
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#1e293b_1px,transparent_1px),linear-gradient(to_bottom,#1e293b_1px,transparent_1px)] bg-[size:4rem_4rem] [mask-image:radial-gradient(ellipse_80%_50%_at_50%_0%,#000,transparent)]"></div>

      <div className="absolute top-0 left-1/4 w-96 h-96 bg-blue-500/20 rounded-full filter blur-[128px] animate-pulse"></div>
      <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-violet-500/20 rounded-full filter blur-[128px] animate-pulse" style={{ animationDelay: '1s' }}></div>

      <main className="relative z-10 max-w-7xl mx-auto px-6 py-12">
        <section className="mb-16">
          <div className="relative group">
            <div className="absolute -inset-1 bg-gradient-to-r from-blue-600 to-violet-600 rounded-3xl blur opacity-25 group-hover:opacity-40 transition duration-1000"></div>
            <div className="relative bg-gradient-to-br from-slate-800/90 to-slate-900/90 backdrop-blur-xl rounded-3xl p-10 border border-slate-700/50 shadow-2xl">
              <h2 className="text-4xl font-bold bg-gradient-to-r from-white to-slate-300 bg-clip-text text-transparent mb-3">
                Welcome back, {user?.username}!
              </h2>
              <p className="text-slate-400 text-lg mb-8">
                Ready to streamline your development workflow today?
              </p>

              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                {quickActions.map((action, index) => {
                  const Icon = action.icon;
                  return (
                    <button
                      key={index}
                      className="group/btn relative overflow-hidden bg-slate-800/50 backdrop-blur-sm border border-slate-700 rounded-xl p-6 hover:border-slate-600 transition-all duration-300 hover:scale-105 hover:shadow-lg hover:shadow-blue-500/10"
                    >
                      <div className={`absolute inset-0 bg-gradient-to-br ${action.gradient} opacity-0 group-hover/btn:opacity-10 transition-opacity duration-300`}></div>
                      <div className="relative flex flex-col items-center space-y-3">
                        <div className={`p-3 rounded-xl bg-gradient-to-br ${action.gradient} shadow-lg`}>
                          <Icon className="w-6 h-6 text-white" />
                        </div>
                        <span className="text-slate-300 font-medium group-hover/btn:text-white transition-colors">
                          {action.label}
                        </span>
                      </div>
                    </button>
                  );
                })}
              </div>
            </div>
          </div>
        </section>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
          <section className="group relative">
            <div className="absolute -inset-0.5 bg-gradient-to-r from-blue-600 to-cyan-600 rounded-2xl blur opacity-0 group-hover:opacity-20 transition duration-500"></div>
            <div className="relative bg-slate-900/80 backdrop-blur-xl rounded-2xl p-8 border border-slate-800 shadow-xl">
              <div className="flex items-center justify-between mb-8">
                <h3 className="text-2xl font-bold text-white flex items-center space-x-3">
                  <div className="p-2 bg-blue-500/10 rounded-lg">
                    <Activity className="w-5 h-5 text-blue-400" />
                  </div>
                  <span>Recent Boards</span>
                </h3>
                <button className="flex items-center space-x-1 text-blue-400 hover:text-blue-300 text-sm font-medium group/link transition-colors">
                  <span>View All</span>
                  <ChevronRight className="w-4 h-4 group-hover/link:translate-x-1 transition-transform" />
                </button>
              </div>

              <div className="space-y-3">
                {recentBoards.map((board) => (
                  <div
                    key={board.id}
                    className="group/item relative overflow-hidden p-5 bg-slate-800/50 backdrop-blur-sm rounded-xl border border-slate-700 hover:border-slate-600 transition-all duration-300 cursor-pointer hover:scale-[1.02] hover:shadow-lg"
                  >
                    <div className={`absolute inset-0 bg-gradient-to-br ${board.color} opacity-0 group-hover/item:opacity-5 transition-opacity duration-300`}></div>
                    <div className="relative flex items-center justify-between">
                      <div className="flex items-center space-x-4">
                        <div className={`w-10 h-10 bg-gradient-to-br ${board.color} rounded-lg shadow-lg flex items-center justify-center`}>
                          <Activity className="w-5 h-5 text-white" />
                        </div>
                        <div>
                          <h4 className="text-slate-200 font-semibold group-hover/item:text-white transition-colors">
                            {board.name}
                          </h4>
                          <div className="flex items-center space-x-2 text-sm text-slate-500 mt-1">
                            <Clock className="w-3 h-3" />
                            <span>{board.updated}</span>
                          </div>
                        </div>
                      </div>
                      <ChevronRight className="w-5 h-5 text-slate-600 group-hover/item:text-slate-400 group-hover/item:translate-x-1 transition-all" />
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </section>

          <section className="group relative">
            <div className="absolute -inset-0.5 bg-gradient-to-r from-violet-600 to-purple-600 rounded-2xl blur opacity-0 group-hover:opacity-20 transition duration-500"></div>
            <div className="relative bg-slate-900/80 backdrop-blur-xl rounded-2xl p-8 border border-slate-800 shadow-xl">
              <div className="flex items-center justify-between mb-8">
                <h3 className="text-2xl font-bold text-white flex items-center space-x-3">
                  <div className="p-2 bg-violet-500/10 rounded-lg">
                    <GitBranch className="w-5 h-5 text-violet-400" />
                  </div>
                  <span>Recent Repositories</span>
                </h3>
                <button className="flex items-center space-x-1 text-violet-400 hover:text-violet-300 text-sm font-medium group/link transition-colors">
                  <span>View All</span>
                  <ChevronRight className="w-4 h-4 group-hover/link:translate-x-1 transition-transform" />
                </button>
              </div>

              <div className="space-y-3">
                {recentRepos.map((repo) => (
                  <div
                    key={repo.id}
                    className="group/item relative overflow-hidden p-5 bg-slate-800/50 backdrop-blur-sm rounded-xl border border-slate-700 hover:border-slate-600 transition-all duration-300 cursor-pointer hover:scale-[1.02] hover:shadow-lg"
                  >
                    <div className="absolute inset-0 bg-gradient-to-br from-violet-500 to-purple-600 opacity-0 group-hover/item:opacity-5 transition-opacity duration-300"></div>
                    <div className="relative flex items-center justify-between">
                      <div className="flex-1">
                        <h4 className="text-slate-200 font-semibold group-hover/item:text-white transition-colors mb-2">
                          {repo.name}
                        </h4>
                        <div className="flex items-center space-x-4 text-sm text-slate-500">
                          <div className="flex items-center space-x-2">
                            <div className="w-3 h-3 bg-blue-400 rounded-full"></div>
                            <span>{repo.language}</span>
                          </div>
                          <div className="flex items-center space-x-1">
                            <Star className="w-3 h-3 fill-amber-400 text-amber-400" />
                            <span>{repo.stars}</span>
                          </div>
                        </div>
                      </div>
                      <button className="px-4 py-2 bg-slate-700 hover:bg-slate-600 rounded-lg text-slate-300 hover:text-white text-sm font-medium transition-all duration-200 hover:scale-105">
                        View
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </section>
        </div>

        <section className="group relative">
          <div className="absolute -inset-0.5 bg-gradient-to-r from-emerald-600 to-teal-600 rounded-2xl blur opacity-0 group-hover:opacity-20 transition duration-500"></div>
          <div className="relative bg-slate-900/80 backdrop-blur-xl rounded-2xl p-8 border border-slate-800 shadow-xl">
            <h3 className="text-2xl font-bold text-white mb-8 flex items-center space-x-3">
              <div className="p-2 bg-emerald-500/10 rounded-lg">
                <Activity className="w-5 h-5 text-emerald-400" />
              </div>
              <span>Recent Activity</span>
            </h3>

            <div className="space-y-3">
              {[
                {
                  action: "created",
                  type: "board",
                  item: "Q4 Planning",
                  time: "2 hours ago",
                  gradient: "from-blue-500 to-cyan-500",
                  icon: Activity,
                },
                {
                  action: "commented",
                  type: "issue",
                  item: "Fix auth bug",
                  time: "4 hours ago",
                  gradient: "from-orange-500 to-red-500",
                  icon: Bug,
                },
                {
                  action: "merged",
                  type: "PR",
                  item: "feature/new-dashboard",
                  time: "1 day ago",
                  gradient: "from-violet-500 to-purple-500",
                  icon: GitBranch,
                },
              ].map((activity, index) => {
                const Icon = activity.icon;
                return (
                  <div
                    key={index}
                    className="group/activity flex items-center space-x-5 p-5 rounded-xl bg-slate-800/50 backdrop-blur-sm border border-slate-700 hover:border-slate-600 transition-all duration-300 hover:scale-[1.01] cursor-pointer"
                  >
                    <div className={`w-12 h-12 bg-gradient-to-br ${activity.gradient} rounded-xl shadow-lg flex items-center justify-center flex-shrink-0`}>
                      <Icon className="w-6 h-6 text-white" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-slate-300 group-hover/activity:text-white transition-colors">
                        You{" "}
                        <span className="text-emerald-400 font-medium">
                          {activity.action}
                        </span>{" "}
                        <span className="text-white font-semibold">
                          {activity.item}
                        </span>
                      </p>
                    </div>
                    <div className="flex items-center space-x-2 text-slate-500 text-sm flex-shrink-0">
                      <Clock className="w-4 h-4" />
                      <span>{activity.time}</span>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </section>
      </main>

      <footer className="relative z-10 border-t border-slate-800/50 mt-16 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto px-6 py-8">
          <div className="flex flex-col sm:flex-row items-center justify-between text-slate-500 text-sm gap-4">
            <div className="flex items-center space-x-4">
              <span className="font-semibold text-slate-400">CodeBoard v1.0</span>
              <span className="hidden sm:inline">•</span>
              <span>Microservices Architecture</span>
            </div>
            <div className="flex items-center space-x-6">
              <button className="hover:text-slate-300 transition-colors duration-200">
                Documentation
              </button>
              <button className="hover:text-slate-300 transition-colors duration-200">
                Support
              </button>
              <button className="hover:text-slate-300 transition-colors duration-200">
                Status
              </button>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
