import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import LoginFormCSS from "./LoginForm.module.scss"
import { AuthService } from 'features/auth/services/AuthService'
import { useState } from 'react'
import { useNavigate } from "react-router-dom";

const LoginForm = (props: { twoFactor: any }) => {

  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const emailChange = (event: any) => {
    setEmail(event.target.value);
  };

  const passwordChange = (event: any) => {
    setPassword(event.target.value);
  };

  const onClick = () => {
    (async function () {
      try {
        const jwt = await AuthService.login({ Email: email, Password: password });
        props.twoFactor();
        // TODO change flow bellow
        localStorage.setItem("token", jwt)
        // navigate("/")

      } catch (error: any) {
        alert(error.response.data);
      }
    })()
  }

  return (
    <div className={LoginFormCSS.form}>
      <InputField className={LoginFormCSS.input} usage="Email" value={email} onChange={emailChange} />
      <InputField className={LoginFormCSS.input} usage="Password" value={password} onChange={passwordChange} />
      <div className={LoginFormCSS.button}>
        <a href="#" className={LoginFormCSS.forgotPassword}>
          Forgot password ?
        </a>
        <span className="alignRight">
          <Button onClick={onClick} text="Sign in" />
        </span>
      </div>
    </div >
  )

}

export default LoginForm