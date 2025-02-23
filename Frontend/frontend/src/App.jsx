import './App.css';
import { BrowserRouter, Routes, Route } from "react-router-dom";
import RegisterPage from './Pages/RegisterPage';
import Login from './Pages/LoginPage';
import OwnerDashboard from './Pages/OwnerDashboard';
function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/login" element={<Login />} />
        <Route path="/owner-dashboard" element={<OwnerDashboard />} />
        <Route />
      </Routes>
    </BrowserRouter >
  );
}

export default App;
