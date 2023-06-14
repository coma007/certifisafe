import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import LoginFormCSS from "./LoginForm.module.scss"
import { AuthService } from 'features/auth/services/AuthService'
import { useState } from 'react'
import { useNavigate } from "react-router-dom";
import * as yup from 'yup' 
import { Formik, Form, Field, ErrorMessage } from 'formik';

const LoginForm = (props: { twoFactor: any }) => {

  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

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

  const schema = yup.object().shape({
    email: yup.string().email().required(),
    password: yup.string().min(8, "password is too short")
    .matches( /[a-z]+/, "needs to contain lowercase letter")
    .matches( /[A-Z]+/, "needs to contain uppercase letter")
    .matches( /[0-9]+/, "needs to contain number letter")
    .required(),
  })


  return (
    <Formik
       initialValues={{
         password: '',
         email: '',
       }}
       validationSchema={schema}
       onSubmit={values => {

       }}
     >
       {({ errors, touched, setFieldValue }) => (
          <Form className={LoginFormCSS.form}>
            <Field name="email" component={InputField} className={LoginFormCSS.input} usage="Email" value={email} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                    setEmail(e.target.value);
                    setFieldValue("email", e.target.value);
                  }}/>
            <ErrorMessage name="email" />

            <Field name="password" component={InputField} className={LoginFormCSS.input} usage="Password" value={password} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                    setPassword(e.target.value);
                    setFieldValue("password", e.target.value);
                  }} />
            <ErrorMessage name="password" />
            <div className={LoginFormCSS.button}>
              <a href="#" className={LoginFormCSS.forgotPassword}>
                Forgot password ?
              </a>
              <span className="alignRight">
                <Button onClick={onClick} text="Sign in" submit={"submit"} />
              </span>
            </div>
          </Form >
      )}
    </Formik>
  )

}

export default LoginForm