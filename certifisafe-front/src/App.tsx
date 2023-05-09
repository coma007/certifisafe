import AppCSS from './App.module.scss';
import HomePage from 'pages/home/HomePage';
import PreviewPage from 'features/certificate/Preview/pages/PreviewPage/PreviewPage';

function App() {
  return (
    <div className={AppCSS.main}>
      <PreviewPage></PreviewPage>
    </div >
  )
}

export default App;
