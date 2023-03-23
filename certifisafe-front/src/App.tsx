import AppCSS from './App.module.scss';
import HomePage from 'pages/home/HomePage';

function App() {
  return (
    <div className={AppCSS.main}>
      <HomePage></HomePage>
    </div >
  )
}

export default App;
