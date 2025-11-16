import { RouterProvider } from 'react-router-dom'
import { QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { Toaster } from 'react-hot-toast'
import { ErrorBoundary } from '@/shared/components/common/ErrorBoundary'
import { Loading } from '@/shared/components/common/Loading'
import { useAuthInit } from '@/shared/hooks/useAuthInit'
import { queryClient } from '@/shared/lib/queryClient'
import { router } from './router'

// ⭐ Extract App content to separate component
const AppContent = () => {
  const { isLoading } = useAuthInit()

  if (isLoading) {
    return <Loading />
  }

  return <RouterProvider router={router} />
}

export const App = () => {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <AppContent />
        <Toaster
          position="top-right"
          toastOptions={{
            duration: 4000,
            style: {
              background: '#363636',
              color: '#fff',
              padding: '16px',
              fontSize: '14px',
            },
            success: {
              duration: 3000,
              iconTheme: {
                primary: '#10b981',
                secondary: '#fff',
              },
            },
            error: {
              duration: 4000,
              iconTheme: {
                primary: '#ef4444',
                secondary: '#fff',
              },
              style: {
                background: '#ef4444',
                color: '#fff',
              },
            },
          }}
        />
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </ErrorBoundary>
  )
}
