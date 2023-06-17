import AppCSS from './App.module.scss';
import Router from 'routes/Router';

function App() {
  
  return (
    <div className={AppCSS.main}>
      <Router></Router>
      <script src="https://www.google.com/recaptcha/api.js" async defer />
    </div >
  )
}

export default App;
