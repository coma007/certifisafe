import LoginPage from 'features/auth/pages/LoginPage/LoginPage';
import RegisterPage from 'features/auth/pages/RegisterPage/RegisterPage';
import CertificateOreviewPage from 'features/certificate/pages/Overview/CertificateOverviewPage';
import RequestOreviewPage from 'features/request/pages/Overview/RequestOverviewPage';
import HomePage from 'pages/home/HomePage';
import React from 'react'
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Guard } from './GuardedRoute';


const Router = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route index element={<HomePage />} />
                <Route path="/login" element={<LoginPage />} />
                <Route path="/register" element={<RegisterPage />} />
                <Route element={<Guard />}>
                    <Route path="/certificates" element={<CertificateOreviewPage />} />
                    <Route path="/requests" element={<RequestOreviewPage />} />
                </Route>
            </Routes>
        </BrowserRouter>
    )
}

export default Router