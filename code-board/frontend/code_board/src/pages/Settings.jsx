import { User, Bell, Shield, Palette, Github } from 'lucide-react';
import Card from '../components/Card';
import Button from '../components/Button';
import { useState } from 'react';

const menuItems = [
  { id: 'profile', label: 'Profile', icon: User },
  { id: 'preferences', label: 'Preferences', icon: Palette },
  { id: 'security', label: 'Security', icon: Shield },
];

export default function Settings() {
  const [activeSection, setActiveSection] = useState('profile');
  const [isDarkMode, setIsDarkMode] = useState(true);

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-100 mb-2">Settings</h1>
        <p className="text-gray-400">Manage your account settings and preferences</p>
      </div>

      <div className="grid grid-cols-12 gap-6">
        <div className="col-span-3">
          <Card>
            <nav className="space-y-2">
              {menuItems.map((item) => {
                const Icon = item.icon;
                return (
                  <button
                    key={item.id}
                    onClick={() => setActiveSection(item.id)}
                    className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${
                      activeSection === item.id
                        ? 'bg-gradient-to-r from-[#38bdf8]/20 to-[#a855f7]/20 border border-[#38bdf8]/30 text-[#38bdf8]'
                        : 'text-gray-400 hover:bg-[#0f172a] hover:text-gray-200'
                    }`}
                  >
                    <Icon className="w-5 h-5" />
                    <span className="font-medium">{item.label}</span>
                  </button>
                );
              })}
            </nav>
          </Card>
        </div>

        <div className="col-span-9 space-y-6">
          {activeSection === 'profile' && (
            <>
              <Card>
                <h2 className="text-xl font-bold text-gray-100 mb-6">Profile Information</h2>
                <div className="flex items-center gap-6 mb-6">
                  <div className="relative">
                    <div className="w-24 h-24 bg-gradient-to-br from-[#38bdf8] to-[#a855f7] rounded-2xl flex items-center justify-center text-white text-3xl font-bold">
                      AD
                    </div>
                    <button className="absolute bottom-0 right-0 w-8 h-8 bg-[#38bdf8] rounded-lg flex items-center justify-center hover:bg-[#a855f7] transition-colors">
                      <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                      </svg>
                    </button>
                  </div>
                  <div>
                    <h3 className="text-lg font-bold text-gray-100 mb-1">Alex Developer</h3>
                    <p className="text-gray-400 mb-2">@alexdev</p>
                    <Button variant="secondary" className="text-sm">Change Avatar</Button>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-2">Full Name</label>
                    <input
                      type="text"
                      defaultValue="Alex Developer"
                      className="w-full px-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-2">Username</label>
                    <input
                      type="text"
                      defaultValue="alexdev"
                      className="w-full px-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                    />
                  </div>
                  <div className="col-span-2">
                    <label className="block text-sm font-medium text-gray-300 mb-2">Email</label>
                    <input
                      type="email"
                      defaultValue="alex@example.com"
                      className="w-full px-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                    />
                  </div>
                  <div className="col-span-2">
                    <label className="block text-sm font-medium text-gray-300 mb-2">Bio</label>
                    <textarea
                      rows={3}
                      defaultValue="Full-stack developer passionate about building great products"
                      className="w-full px-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors resize-none"
                    />
                  </div>
                  <div className="col-span-2">
                    <label className="block text-sm font-medium text-gray-300 mb-2">GitHub Profile</label>
                    <div className="relative">
                      <Github className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                      <input
                        type="text"
                        defaultValue="github.com/alexdev"
                        className="w-full pl-10 pr-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                      />
                    </div>
                  </div>
                </div>

                <div className="flex justify-end gap-3 mt-6">
                  <Button variant="secondary">Cancel</Button>
                  <Button variant="primary">Save Changes</Button>
                </div>
              </Card>
            </>
          )}

          {activeSection === 'preferences' && (
            <>
              <Card>
                <h2 className="text-xl font-bold text-gray-100 mb-6">Appearance</h2>
                <div className="space-y-4">
                  <div className="flex items-center justify-between py-3 border-b border-gray-800">
                    <div>
                      <div className="font-medium text-gray-100 mb-1">Dark Mode</div>
                      <div className="text-sm text-gray-400">Enable dark theme across the platform</div>
                    </div>
                    <button
                      onClick={() => setIsDarkMode(!isDarkMode)}
                      className={`relative w-12 h-6 rounded-full transition-colors ${
                        isDarkMode ? 'bg-[#38bdf8]' : 'bg-gray-700'
                      }`}
                    >
                      <div
                        className={`absolute top-1 left-1 w-4 h-4 bg-white rounded-full transition-transform ${
                          isDarkMode ? 'translate-x-6' : ''
                        }`}
                      />
                    </button>
                  </div>
                </div>
              </Card>

              <Card>
                <h2 className="text-xl font-bold text-gray-100 mb-6">Notifications</h2>
                <div className="space-y-4">
                  {[
                    { label: 'Email Notifications', desc: 'Receive email updates about your activity' },
                    { label: 'Push Notifications', desc: 'Receive push notifications on this device' },
                    { label: 'Project Updates', desc: 'Get notified about project changes' },
                    { label: 'Team Mentions', desc: 'Receive notifications when someone mentions you' },
                  ].map((item, idx) => (
                    <div key={idx} className="flex items-center justify-between py-3 border-b border-gray-800 last:border-0">
                      <div>
                        <div className="font-medium text-gray-100 mb-1">{item.label}</div>
                        <div className="text-sm text-gray-400">{item.desc}</div>
                      </div>
                      <button className="relative w-12 h-6 rounded-full bg-[#38bdf8]">
                        <div className="absolute top-1 right-1 w-4 h-4 bg-white rounded-full" />
                      </button>
                    </div>
                  ))}
                </div>
              </Card>
            </>
          )}

          {activeSection === 'security' && (
            <>
              <Card>
                <h2 className="text-xl font-bold text-gray-100 mb-6">Change Password</h2>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-2">Current Password</label>
                    <input
                      type="password"
                      className="w-full px-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-2">New Password</label>
                    <input
                      type="password"
                      className="w-full px-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-2">Confirm New Password</label>
                    <input
                      type="password"
                      className="w-full px-4 py-3 bg-[#0f172a] border border-gray-700 rounded-xl text-gray-100 focus:outline-none focus:border-[#38bdf8] transition-colors"
                    />
                  </div>
                </div>
                <div className="flex justify-end mt-6">
                  <Button variant="primary">Update Password</Button>
                </div>
              </Card>

              <Card>
                <h2 className="text-xl font-bold text-gray-100 mb-6">Two-Factor Authentication</h2>
                <div className="flex items-start gap-4 mb-6">
                  <Shield className="w-12 h-12 text-[#10b981] mt-1" />
                  <div className="flex-1">
                    <h3 className="font-medium text-gray-100 mb-2">Secure Your Account</h3>
                    <p className="text-sm text-gray-400 mb-4">
                      Two-factor authentication adds an extra layer of security to your account by requiring a code in addition to your password.
                    </p>
                    <Button variant="primary">Enable 2FA</Button>
                  </div>
                </div>
              </Card>

              <Card>
                <h2 className="text-xl font-bold text-gray-100 mb-6 text-red-400">Danger Zone</h2>
                <div className="space-y-4">
                  <div className="flex items-center justify-between py-3 border-b border-gray-800">
                    <div>
                      <div className="font-medium text-gray-100 mb-1">Delete Account</div>
                      <div className="text-sm text-gray-400">Permanently delete your account and all data</div>
                    </div>
                    <Button variant="secondary" className="border-red-500 text-red-400 hover:bg-red-500/10">
                      Delete
                    </Button>
                  </div>
                </div>
              </Card>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
