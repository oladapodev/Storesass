import { Routes, Route } from 'react-router-dom'
import { Layout } from '@/components/Layout'
import { LandingPage } from '@/pages/LandingPage'
import { StoresPage } from '@/pages/StoresPage'
import { StoreDetailPage } from '@/pages/StoreDetailPage'
import { DashboardPage } from '@/pages/DashboardPage'

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<LandingPage />} />
        <Route path="/stores" element={<StoresPage />} />
        <Route path="/stores/:slug" element={<StoreDetailPage />} />
        <Route path="/dashboard" element={<DashboardPage />} />
      </Route>
    </Routes>
  )
}
