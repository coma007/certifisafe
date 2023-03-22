import AppCSS from './App.module.scss';
import LoginPage from 'features/auth/Login/pages/LoginPage/LoginPage';
import RegisterPage from 'features/auth/Register/pages/RegisterPage/RegisterPage';
import HomePage from 'pages/Home/HomePage';

function App() {
  return (
    <div className={AppCSS.main}>
      <HomePage></HomePage>
    </div >
  )
}

export default App;
