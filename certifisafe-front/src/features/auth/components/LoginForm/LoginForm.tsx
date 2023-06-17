import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import LoginFormCSS from "./LoginForm.module.scss"
import { AuthService } from 'features/auth/services/AuthService'
import { useRef, useState } from 'react'
import { useNavigate } from "react-router-dom";
import * as yup from 'yup' 
import { Formik, Form, Field, ErrorMessage } from 'formik';
// import reCAPTCHA from "react-google-recaptcha"
import ReCAPTCHA from 'react-google-recaptcha'
import ReactDOM, { render } from 'react-dom'
import ErrorMsg from 'components/error/ErrorMsg'

const LoginForm = (props: { twoFactor: any, resetPassword : any }) => {

  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const onClick = () => {
    (async function () {
      try {
          const jwt = await AuthService.login({ Email: email, Password: password });
          props.twoFactor();
          
          // TODO change flow bellow
          
      } catch (error: any) {
        if (error.response.status == 403) {
          props.resetPassword();
        } else {
          alert(error.response.data);
        }
    }
    })()
  }

  const schema = yup.object().shape({
    email: yup.string().email().required(),
    password: yup.string().min(8, "password is too short")
    .matches( /[a-z]+/, "needs to contain lowercase letter")
    .matches( /[A-Z]+/, "needs to contain uppercase letter")
    .matches( /[0-9]+/, "needs to contain number letter")
    .required(),
  })

  const captchaRef: any = useRef(null)

  const handleSubmit = async () => {
    try {
      const token = captchaRef.current?.getValue();
      captchaRef.current.reset();
      const jwt = await AuthService.login({ Email: email, Password: password, Token: token});
      props.twoFactor();
      // TODO change flow bellow
      localStorage.setItem("token", jwt)
      // navigate("/")
  } catch (error: any) {
    alert(error.response.data);
}
}


  return ( 
    <Formik
       initialValues={{
         password: '',
         email: '',
         token: '',
       }}
       validationSchema={schema}
       validateOnChange
       onSubmit={handleSubmit}
     >
       {({ errors, touched, setFieldValue, validateForm, isValid, handleSubmit }) => (

          <Form className={LoginFormCSS.form}> 
            <Field name="email" component={InputField} className={LoginFormCSS.input} usage="Email" value={email} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                    setEmail(e.target.value);
                    setFieldValue("email", e.target.value);
                  }}/>
            <ErrorMsg val={errors["email"]} />

            <Field name="password" component={InputField} className={LoginFormCSS.input} usage="Password" value={password} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                    setPassword(e.target.value);
                    setFieldValue("password", e.target.value);
                  }} />
            <ErrorMsg val={errors["password"]} />
            <div className={LoginFormCSS.button}>
              <a href="/passwordRecovery" className={LoginFormCSS.forgotPassword}>
                Forgot password ?
              </a>
              <span className="alignRight">
                <Button onClick={null} text="Sign in" submit={"submit"} />
              </span>
            </div>
            <ReCAPTCHA className='center' sitekey={process.env.REACT_APP_SITE_KEY as string}  ref={captchaRef}/>
          </Form >
      )}
    </Formik>      
    )

}


export default LoginForm