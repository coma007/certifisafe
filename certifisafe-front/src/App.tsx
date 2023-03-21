import AppCSS from './App.module.scss';
import LoginPage from 'features/auth/Login/pages/LoginPage/LoginPage';
import RegisterPage from 'features/auth/Register/pages/RegisterPage/RegisterPage';

function App() {
  return (
    <div className={AppCSS.main}>
      <LoginPage></LoginPage>
    </div >
  )
}

export default App;
