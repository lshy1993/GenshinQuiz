import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { SWRConfig } from 'swr'
import fetcher from './api/fetcher/fetcher.ts'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <SWRConfig 
      value={{
        fetcher,
        refreshInterval: 30000, // 30秒自动刷新
        revalidateOnFocus: true,
        revalidateOnReconnect: true,
        errorRetryCount: 3,
      }}
    > <App />
    </SWRConfig>
  </StrictMode>,
)
