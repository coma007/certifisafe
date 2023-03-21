import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import React from 'react'
import LoginFormCSS from "./LoginForm.module.scss"

const LoginForm = () => {
  return (
    <div className={LoginFormCSS.form}>
      <InputField className={LoginFormCSS.input} usage="Email" />
      <InputField className={LoginFormCSS.input} usage="Password" />
      <div className={LoginFormCSS.button}>
        <a href="#" className={LoginFormCSS.forgotPassword}>Forgot password ?</a>
        <span className={LoginFormCSS.alignRight}>
          <Button text="Sign in"></Button>
        </span>
      </div>
    </div >
  )
}

export default LoginForm