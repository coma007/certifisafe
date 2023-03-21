import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import RegisterFormCSS from './RegisterForm.module.scss'

const RegisterForm = () => {
  return (
    <div className={RegisterFormCSS.form}>
      <InputField className={RegisterFormCSS.inlineInput} usage="First name" />
      <InputField className={`alignRight ${RegisterFormCSS.inlineInput}`} usage="Last name" />
      <InputField className={RegisterFormCSS.input} usage="Email" />
      <InputField className={RegisterFormCSS.input} usage="Password" />
      <div className={RegisterFormCSS.button}>
        <span className="alignRight">
          <Button text="Get started" />
        </span>
      </div>
    </div >
  )
}

export default RegisterForm