import AppCSS from './App.module.scss';
import Router from 'routes/Router';

function App() {
  
  return (
    <div className={AppCSS.main}>
      <Router></Router>
    </div >
  )
}

export default App;
