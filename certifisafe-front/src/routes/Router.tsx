import LoginPage, { Logout } from 'features/auth/pages/LoginPage/LoginPage';
import RegisterPage from 'features/auth/pages/RegisterPage/RegisterPage';
import CertificateOreviewPage from 'features/certificate/pages/Overview/CertificateOverviewPage';
import CertificateVerifyPage from 'features/certificate/pages/Verify/CertificateVerifyPage';
import RequestOreviewPage from 'features/request/pages/Overview/RequestOverviewPage';
import HomePage from 'pages/home/HomePage';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthGuard, NonAuthGuard } from './GuardedRoute';

const Router = () => {

    return (
        <BrowserRouter>
            <Routes>
                <Route index element={<HomePage />} />
                <Route element={<NonAuthGuard />}>
                    <Route path="/login" element={<LoginPage />} />
                    <Route path="/register" element={<RegisterPage />} />
                </Route>
                <Route element={<AuthGuard />}>
                    <Route path="/certificates" element={<CertificateOreviewPage />} />
                    <Route path="/verify" element={<CertificateVerifyPage />} />
                    <Route path="/requests" element={<RequestOreviewPage />} />
                    <Route path="/logout" element={<Logout />} />
                </Route>
            </Routes>
        </BrowserRouter>
    );
};

export default Router;