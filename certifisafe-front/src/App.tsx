import AppCSS from './App.module.scss';
import { BrowserRouter, Routes, Route } from "react-router-dom";
import RequestOreviewPage from 'features/request/pages/Overview/RequestOverviewPage';
import CertificateOreviewPage from 'features/certificate/pages/Overview/CertificateOverviewPage';
import LoginPage from 'features/auth/pages/LoginPage/LoginPage';
import HomePage from 'pages/home/HomePage';
import RegisterPage from 'features/auth/pages/RegisterPage/RegisterPage';

function App() {
  return (
    <div className={AppCSS.main}>
      <BrowserRouter>
        <Routes>
          <Route index element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/certificates" element={<CertificateOreviewPage />} />
          <Route path="/requests" element={<RequestOreviewPage />} />
        </Routes>
      </BrowserRouter>
    </div >
  )
}

export default App;
