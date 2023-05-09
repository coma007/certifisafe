import AppCSS from './App.module.scss';
import CertificateOreviewPage from 'features/certificate/Preview/pages/Overview/CertificateOverviewPage';

function App() {
  return (
    <div className={AppCSS.main}>
      <CertificateOreviewPage></CertificateOreviewPage>
    </div >
  )
}

export default App;
