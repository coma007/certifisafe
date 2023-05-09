import AppCSS from './App.module.scss';
import RequestOreviewPage from 'features/request/pages/Overview/RequestOverviewPage';
import CertificateOreviewPage from 'features/certificate/pages/Overview/CertificateOverviewPage';
import LoginPage from 'features/auth/pages/LoginPage/LoginPage';

function App() {
  return (
    <div className={AppCSS.main}>
      <CertificateOreviewPage></CertificateOreviewPage>
    </div >
  )
}

export default App;
