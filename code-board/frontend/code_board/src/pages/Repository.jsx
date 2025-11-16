import { Star, GitFork, Eye, Code2, GitCommit, AlertCircle, GitPullRequest, Settings } from 'lucide-react';
import Button from '../components/Button';
import { useState } from 'react';

const tabs = [
  { id: 'code', label: 'Code', icon: Code2 },
  { id: 'commits', label: 'Commits', icon: GitCommit },
  { id: 'issues', label: 'Issues', icon: AlertCircle, count: 8 },
  { id: 'pulls', label: 'Pull Requests', icon: GitPullRequest, count: 3 },
  { id: 'settings', label: 'Settings', icon: Settings },
];

const files = [
  { name: 'src', type: 'folder', updated: '2 days ago' },
  { name: 'public', type: 'folder', updated: '1 week ago' },
  { name: 'node_modules', type: 'folder', updated: '3 days ago' },
  { name: 'package.json', type: 'file', updated: '2 days ago' },
  { name: 'README.md', type: 'file', updated: '1 week ago' },
  { name: 'vite.config.js', type: 'file', updated: '2 weeks ago' },
  { name: '.gitignore', type: 'file', updated: '3 weeks ago' },
];

export default function Repository() {
  const [activeTab, setActiveTab] = useState('code');

  return (
    <div>
      <div className="bg-[#1e293b] border border-gray-800 rounded-2xl p-6 mb-6">
        <div className="flex items-start justify-between mb-4">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-3xl font-bold text-gray-100">demo-project</h1>
              <span className="px-3 py-1 bg-[#0f172a] border border-[#10b981] text-[#10b981] rounded-lg text-xs font-medium">
                Public
              </span>
            </div>
            <p className="text-gray-400">A modern web application built with React and TailwindCSS</p>
          </div>
          <div className="flex items-center gap-3">
            <Button variant="secondary" className="flex items-center gap-2">
              <Eye className="w-4 h-4" />
              Watch
              <span className="ml-1 px-2 py-0.5 bg-[#0f172a] rounded-full text-xs">12</span>
            </Button>
            <Button variant="secondary" className="flex items-center gap-2">
              <Star className="w-4 h-4" />
              Star
              <span className="ml-1 px-2 py-0.5 bg-[#0f172a] rounded-full text-xs">42</span>
            </Button>
            <Button variant="secondary" className="flex items-center gap-2">
              <GitFork className="w-4 h-4" />
              Fork
              <span className="ml-1 px-2 py-0.5 bg-[#0f172a] rounded-full text-xs">8</span>
            </Button>
          </div>
        </div>

        <div className="flex items-center gap-6 text-sm text-gray-400">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-[#38bdf8] rounded-full"></div>
            JavaScript
          </div>
          <div className="flex items-center gap-2">
            <GitCommit className="w-4 h-4" />
            156 commits
          </div>
          <div className="flex items-center gap-2">
            <GitFork className="w-4 h-4" />
            main branch
          </div>
          <div>Updated 2 hours ago</div>
        </div>
      </div>

      <div className="border-b border-gray-800 mb-6">
        <div className="flex gap-1">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 px-4 py-3 border-b-2 transition-colors ${
                  activeTab === tab.id
                    ? 'border-[#38bdf8] text-[#38bdf8]'
                    : 'border-transparent text-gray-400 hover:text-gray-200'
                }`}
              >
                <Icon className="w-4 h-4" />
                {tab.label}
                {tab.count && (
                  <span className="px-2 py-0.5 bg-[#1e293b] border border-gray-700 rounded-full text-xs">
                    {tab.count}
                  </span>
                )}
              </button>
            );
          })}
        </div>
      </div>

      {activeTab === 'code' && (
        <div className="grid grid-cols-12 gap-6">
          <div className="col-span-3">
            <div className="bg-[#1e293b] border border-gray-800 rounded-2xl p-4">
              <h3 className="font-bold text-gray-100 mb-4 flex items-center gap-2">
                <Code2 className="w-4 h-4" />
                Files
              </h3>
              <div className="space-y-1">
                {files.map((file) => (
                  <button
                    key={file.name}
                    className="w-full text-left px-3 py-2 rounded-lg hover:bg-[#0f172a] transition-colors group"
                  >
                    <div className="flex items-center gap-2">
                      {file.type === 'folder' ? (
                        <svg className="w-4 h-4 text-[#38bdf8]" fill="currentColor" viewBox="0 0 20 20">
                          <path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
                        </svg>
                      ) : (
                        <svg className="w-4 h-4 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" />
                        </svg>
                      )}
                      <span className="text-sm text-gray-300 group-hover:text-[#38bdf8] transition-colors">
                        {file.name}
                      </span>
                    </div>
                  </button>
                ))}
              </div>
            </div>
          </div>

          <div className="col-span-9">
            <div className="bg-[#1e293b] border border-gray-800 rounded-2xl overflow-hidden">
              <div className="px-4 py-3 border-b border-gray-800 flex items-center justify-between">
                <span className="text-sm text-gray-400 font-mono">src/App.jsx</span>
                <Button variant="ghost" className="text-xs">
                  Copy
                </Button>
              </div>
              <div className="p-4 font-mono text-sm overflow-x-auto">
                <pre className="text-gray-300">
                  <code>
                    <span className="text-[#a855f7]">import</span> <span className="text-gray-300">{'{'}</span> <span className="text-[#38bdf8]">useState</span> <span className="text-gray-300">{'}'}</span> <span className="text-[#a855f7]">from</span> <span className="text-[#10b981]">'react'</span><span className="text-gray-300">;</span>{'\n'}
                    <span className="text-[#a855f7]">import</span> <span className="text-[#38bdf8]">Dashboard</span> <span className="text-[#a855f7]">from</span> <span className="text-[#10b981]">'./components/Dashboard'</span><span className="text-gray-300">;</span>{'\n'}
                    {'\n'}
                    <span className="text-[#a855f7]">function</span> <span className="text-[#f59e0b]">App</span><span className="text-gray-300">() {'{'}</span>{'\n'}
                    {'  '}<span className="text-[#a855f7]">const</span> <span className="text-gray-300">[</span><span className="text-[#38bdf8]">user</span><span className="text-gray-300">,</span> <span className="text-[#38bdf8]">setUser</span><span className="text-gray-300">]</span> <span className="text-[#a855f7]">=</span> <span className="text-[#f59e0b]">useState</span><span className="text-gray-300">(</span><span className="text-[#a855f7]">null</span><span className="text-gray-300">);</span>{'\n'}
                    {'\n'}
                    {'  '}<span className="text-[#a855f7]">return</span> <span className="text-gray-300">(</span>{'\n'}
                    {'    '}<span className="text-gray-300">{'<'}</span><span className="text-[#38bdf8]">div</span> <span className="text-[#a855f7]">className</span><span className="text-gray-300">=</span><span className="text-[#10b981]">"app"</span><span className="text-gray-300">{'>'}</span>{'\n'}
                    {'      '}<span className="text-gray-300">{'<'}</span><span className="text-[#38bdf8]">Dashboard</span> <span className="text-[#a855f7]">user</span><span className="text-gray-300">={'{'}user{'}'}</span> <span className="text-gray-300">/{'>'}</span>{'\n'}
                    {'    '}<span className="text-gray-300">{'</'}</span><span className="text-[#38bdf8]">div</span><span className="text-gray-300">{'>'}</span>{'\n'}
                    {'  '}<span className="text-gray-300">);</span>{'\n'}
                    <span className="text-gray-300">{'}'}</span>{'\n'}
                    {'\n'}
                    <span className="text-[#a855f7]">export default</span> <span className="text-[#38bdf8]">App</span><span className="text-gray-300">;</span>
                  </code>
                </pre>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
