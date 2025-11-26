import { Outlet, Link, useLocation } from 'react-router-dom'

const Layout = () => {
  const location = useLocation()

  const isActive = (path: string) => {
    return location.pathname === path
      ? 'bg-primary-600 text-white'
      : 'text-gray-300 hover:bg-gray-700 hover:text-white'
  }

  return (
    <div className="min-h-screen bg-gray-900">
      {/* Navigation */}
      <nav className="bg-gray-800 border-b border-gray-700">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <h1 className="text-white text-xl font-bold">
                  🔧 Go Formal Verification
                </h1>
              </div>
              <div className="hidden md:block">
                <div className="ml-10 flex items-baseline space-x-4">
                  <Link
                    to="/dashboard"
                    className={`px-3 py-2 rounded-md text-sm font-medium ${isActive(
                      '/dashboard'
                    )}`}
                  >
                    Dashboard
                  </Link>
                  <Link
                    to="/analysis"
                    className={`px-3 py-2 rounded-md text-sm font-medium ${isActive(
                      '/analysis'
                    )}`}
                  >
                    Analysis
                  </Link>
                  <Link
                    to="/patterns"
                    className={`px-3 py-2 rounded-md text-sm font-medium ${isActive(
                      '/patterns'
                    )}`}
                  >
                    Patterns
                  </Link>
                  <Link
                    to="/projects"
                    className={`px-3 py-2 rounded-md text-sm font-medium ${isActive(
                      '/projects'
                    )}`}
                  >
                    Projects
                  </Link>
                </div>
              </div>
            </div>
            <div className="flex items-center">
              <span className="text-sm text-gray-400">v1.0.0</span>
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <Outlet />
      </main>
    </div>
  )
}

export default Layout

