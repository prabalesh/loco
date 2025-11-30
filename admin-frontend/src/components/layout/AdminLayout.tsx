// src/components/layout/AdminLayout.tsx
import { useState } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import {
  LayoutDashboard,
  User,
  LogOut,
  ChevronLeft,
  ChevronRight,
  Users,
  Shield,
} from 'lucide-react'
import { useAuthStore } from '../../features/auth/store/authStore'
import { adminAuthApi } from '../../api/adminApi'
import toast from 'react-hot-toast'

export const AdminLayout = () => {
  const [collapsed, setCollapsed] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()

  const handleLogout = async () => {
    try {
      await adminAuthApi.logout()
      logout()
      navigate('/login')
      toast.success('Logged out successfully')
    } catch {
      toast.error('Logout failed')
    }
  }

  const menuItems = [
    {
      key: '/',
      icon: <LayoutDashboard className="w-5 h-5" />,
      label: 'Dashboard',
      onClick: () => navigate('/'),
    },
    {
      key: '/users',
      icon: <Users className="w-5 h-5" />,
      label: 'Users',
      onClick: () => navigate('/users'),
    },
  ]

  return (
    <div className="min-h-screen flex bg-gray-50">
      {/* Sidebar */}
      <aside
        className={`flex flex-col bg-white border-r border-gray-200 transition-width duration-300 ${
          collapsed ? 'w-20' : 'w-56'
        }`}
      >
        <div className="flex items-center justify-center h-16 border-b border-gray-200">
          <Shield className="text-blue-600 w-7 h-7" />
          {!collapsed && <span className="ml-3 text-lg font-semibold text-gray-700">Admin</span>}
        </div>

        <nav className="flex-1 px-2 py-4 space-y-1">
          {menuItems.map(({ key, icon, label, onClick }) => {
            const isActive = location.pathname === key
            return (
              <button
                key={key}
                onClick={onClick}
                className={`flex items-center w-full px-3 py-2 rounded-md text-left text-gray-700 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                  isActive ? 'bg-blue-100 font-medium' : ''
                }`}
              >
                <span>{icon}</span>
                {!collapsed && <span className="ml-3">{label}</span>}
              </button>
            )
          })}
        </nav>

        <button
          onClick={() => setCollapsed(!collapsed)}
          className="flex items-center justify-center h-12 border-t border-gray-200 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
          aria-label="Toggle Sidebar"
        >
          {collapsed ? <ChevronRight className="w-5 h-5 text-gray-600" /> : <ChevronLeft className="w-5 h-5 text-gray-600" />}
        </button>
      </aside>

      {/* Main content */}
      <div className="flex flex-col flex-1 min-h-screen">
        {/* Header */}
        <header className="flex items-center justify-between h-16 px-6 bg-white border-b border-gray-200">
          <div className="text-lg font-semibold text-gray-900">
            {location.pathname === '/' ? 'Dashboard' : location.pathname.replace('/', '').replace('-', ' ').toUpperCase()}
          </div>
          <button
            className="flex items-center space-x-2 text-gray-700 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded"
            onClick={handleLogout}
            aria-label="Logout"
          >
            <User className="w-5 h-5" />
            <span className="hidden sm:inline">{user?.username || 'Admin'}</span>
            <LogOut className="w-5 h-5" />
          </button>
        </header>

        {/* Content */}
        <main className="flex-grow p-6 bg-gray-50 overflow-auto">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
