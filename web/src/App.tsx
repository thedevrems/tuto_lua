import { Routes, Route, Navigate } from 'react-router-dom'
import HomePage from './pages/HomePage'
import PricingPage from './pages/PricingPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import LearnPage from './pages/LearnPage'

/** Top-level route table. Public marketing + auth pages, plus the learning IDE. */
export default function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/pricing" element={<PricingPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/learn" element={<LearnPage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
