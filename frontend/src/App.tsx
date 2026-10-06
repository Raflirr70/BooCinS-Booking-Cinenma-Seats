import { BrowserRouter, Routes, Route } from 'react-router-dom'
import LoginPage from './pages/auth/login_page'
import RegisterPage from './pages/auth/register_page'
import HomePage from './pages/public/home_page'
import RequireAuth from './components/RequireAuth'
import AdminLayout from './layouts/admin_layout'
import DashboardLayout from './layouts/dashboard_layout'

import AdminDashboardPage from './pages/admin/admin_dashboard_page'
import AdminFilmsPage from './pages/admin/admin_film_page'
import AdminRoomsPage from './pages/admin/admin_room_page'
import AdminSchedulesPage from './pages/admin/admin_schedule_page'
import AdminPromosPage from './pages/admin/admin_promo_page'

import ManagerDashboardPage from './pages/manager/manager_dashboard_page'
import SuperadminDashboardPage from './pages/superadmin/superadmin_dashboard_page'
import StaffCashierPage from './pages/staff/staff_cashier_page'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

        <Route path="/admin" element={<RequireAuth roles={['admin']}><AdminLayout /></RequireAuth>}>
          <Route index element={<AdminDashboardPage />} />
          <Route path="films" element={<AdminFilmsPage />} />
          <Route path="rooms" element={<AdminRoomsPage />} />
          <Route path="schedules" element={<AdminSchedulesPage />} />
          <Route path="promos" element={<AdminPromosPage />} />
        </Route>

        <Route path="/manager" element={<RequireAuth roles={['manager']}><DashboardLayout /></RequireAuth>}>
          <Route index element={<ManagerDashboardPage />} />
        </Route>

        <Route path="/superadmin" element={<RequireAuth roles={['super_admin']}><DashboardLayout /></RequireAuth>}>
          <Route index element={<SuperadminDashboardPage />} />
        </Route>

        <Route path="/staff" element={<RequireAuth roles={['staff']}><DashboardLayout /></RequireAuth>}>
          <Route index element={<StaffCashierPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}