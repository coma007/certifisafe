import Banner from 'components/navigation/Banner/Banner';
import AppCSS from './App.module.scss';
import LoginForm from 'features/auth/Login/components/LoginForm/LoginForm';
import LoginPage from 'features/auth/Login/pages/LoginPage/LoginPage';

function App() {
  return (
    <div className={AppCSS.main}>
      <LoginPage></LoginPage>
    </div >
  )
}

export default App;
