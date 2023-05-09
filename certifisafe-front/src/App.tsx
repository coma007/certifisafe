import AppCSS from './App.module.scss';
import HomePage from 'pages/home/HomePage';
import RequestOreviewPage from 'features/request/Preview/pages/Overview/RequestOverviewPage';
import CertificateOreviewPage from 'features/certificate/Preview/pages/Overview/CertificateOverviewPage';

function App() {
  return (
    <div className={AppCSS.main}>
      <RequestOreviewPage></RequestOreviewPage>
    </div >
  )
}

export default App;
