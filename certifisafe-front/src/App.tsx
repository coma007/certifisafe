import AppCSS from './App.module.scss';
import Router from 'routes/Router';
import { useLocation } from 'react-router-dom';
import { useEffect } from 'react';
import { AuthService } from 'features/auth/services/AuthService';

function App() {
  
  return (
    <div className={AppCSS.main}>
      <Router></Router>
    </div >
  )
}

export default App;
