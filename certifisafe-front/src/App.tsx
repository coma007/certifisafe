import AppCSS from './App.module.scss';
import RequestOreviewPage from 'features/request/pages/Overview/RequestOverviewPage';
import CertificateOreviewPage from 'features/certificate/pages/Overview/CertificateOverviewPage';
import LoginPage from 'features/auth/pages/LoginPage/LoginPage';
import HomePage from 'pages/home/HomePage';
import RegisterPage from 'features/auth/pages/RegisterPage/RegisterPage';
import Router from 'routes/Router';

function App() {
  return (
    <div className={AppCSS.main}>
      <Router></Router>
    </div >
  )
}

export default App;
