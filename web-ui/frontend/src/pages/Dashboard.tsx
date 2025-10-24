import { useEffect, useState } from 'react'
import { apiClient } from '../utils/api'

interface Stats {
  totalProjects: number
  totalAnalyses: number
  totalPatterns: number
  recentActivity: string
}

const Dashboard = () => {
  const [stats, setStats] = useState<Stats>({
    totalProjects: 0,
    totalAnalyses: 0,
    totalPatterns: 30,
    recentActivity: 'Loading...',
  })
  const [health, setHealth] = useState<{ status: string; version: string } | null>(null)

  useEffect(() => {
    // 检查后端健康状态
    apiClient.get('/health').then((res) => {
      setHealth(res.data)
      setStats((prev) => ({
        ...prev,
        recentActivity: `Backend connected (${res.data.version})`,
      }))
    }).catch(() => {
      setStats((prev) => ({
        ...prev,
        recentActivity: 'Backend disconnected',
      }))
    })
  }, [])

  return (
    <div className="px-4 py-6 sm:px-0">
      <h1 className="text-3xl font-bold text-white mb-8">Dashboard</h1>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4 mb-8">
        <div className="bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <span className="text-4xl">📊</span>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-400 truncate">
                    Total Projects
                  </dt>
                  <dd className="text-3xl font-semibold text-white">
                    {stats.totalProjects}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <span className="text-4xl">🔍</span>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-400 truncate">
                    Total Analyses
                  </dt>
                  <dd className="text-3xl font-semibold text-white">
                    {stats.totalAnalyses}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <span className="text-4xl">🎨</span>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-400 truncate">
                    Available Patterns
                  </dt>
                  <dd className="text-3xl font-semibold text-white">
                    {stats.totalPatterns}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <span className="text-4xl">
                  {health?.status === 'ok' ? '✅' : '🔴'}
                </span>
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-400 truncate">
                    Backend Status
                  </dt>
                  <dd className="text-lg font-semibold text-white">
                    {health?.status || 'Checking...'}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Welcome Section */}
      <div className="bg-gray-800 shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-bold text-white mb-4">
          Welcome to Go Formal Verification Framework
        </h2>
        <p className="text-gray-300 mb-4">
          Analyze your Go code for concurrency issues, generate verified patterns,
          and ensure type safety with formal methods.
        </p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
          <div className="border border-gray-700 rounded-lg p-4">
            <h3 className="text-lg font-semibold text-white mb-2">
              🔍 Run Analysis
            </h3>
            <p className="text-gray-400 text-sm mb-3">
              Analyze your code for concurrency issues, deadlocks, and data races.
            </p>
            <button className="w-full bg-primary-600 hover:bg-primary-700 text-white font-medium py-2 px-4 rounded">
              Start Analysis
            </button>
          </div>

          <div className="border border-gray-700 rounded-lg p-4">
            <h3 className="text-lg font-semibold text-white mb-2">
              🎨 Generate Pattern
            </h3>
            <p className="text-gray-400 text-sm mb-3">
              Generate verified concurrency patterns with CSP formal definitions.
            </p>
            <button className="w-full bg-primary-600 hover:bg-primary-700 text-white font-medium py-2 px-4 rounded">
              Browse Patterns
            </button>
          </div>

          <div className="border border-gray-700 rounded-lg p-4">
            <h3 className="text-lg font-semibold text-white mb-2">
              📊 Create Project
            </h3>
            <p className="text-gray-400 text-sm mb-3">
              Manage multiple projects and track analysis history over time.
            </p>
            <button className="w-full bg-primary-600 hover:bg-primary-700 text-white font-medium py-2 px-4 rounded">
              New Project
            </button>
          </div>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="bg-gray-800 shadow rounded-lg p-6">
        <h2 className="text-xl font-bold text-white mb-4">Recent Activity</h2>
        <div className="text-gray-300">
          {stats.recentActivity}
        </div>
      </div>
    </div>
  )
}

export default Dashboard

