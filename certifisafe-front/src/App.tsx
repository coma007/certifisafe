import AppCSS from './App.module.scss';
import CertificateOreviewPage from 'features/certificate/Preview/pages/Overview/CertificateOverviewPage';
import RequestOreviewPage from 'features/request/Preview/pages/Overview/RequestOverviewPage';

function App() {
  return (
    <div className={AppCSS.main}>
      <RequestOreviewPage></RequestOreviewPage>
    </div >
  )
}

export default App;
