import AppCSS from './App.module.scss';
import RequestOreviewPage from 'features/request/pages/Overview/RequestOverviewPage';

function App() {
  return (
    <div className={AppCSS.main}>
      <RequestOreviewPage></RequestOreviewPage>
    </div >
  )
}

export default App;
