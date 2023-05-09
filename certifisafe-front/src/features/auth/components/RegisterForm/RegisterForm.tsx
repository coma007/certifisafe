import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import RegisterFormCSS from './RegisterForm.module.scss'

const RegisterForm = () => {
  return (
    <div className={RegisterFormCSS.form}>
      <InputField className={RegisterFormCSS.inlineInput} usage="First name" value='' onChange={null}/>
      <InputField className={`alignRight ${RegisterFormCSS.inlineInput}`} usage="Last name" value='' onChange={null}/>
      <InputField className={RegisterFormCSS.input} usage="Email" value='' onChange={null}/>
      <InputField className={RegisterFormCSS.input} usage="Password" value='' onChange={null}/>
      <div className={RegisterFormCSS.button}>
        <span className="alignRight">
          <Button onClick={null} text="Get started" />
        </span>
      </div>
    </div >
  )
}

export default RegisterForm