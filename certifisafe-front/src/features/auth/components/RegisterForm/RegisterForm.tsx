import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import RegisterFormCSS from './RegisterForm.module.scss'
import { useState } from 'react';
import { AuthService } from 'features/auth/services/AuthService';
import { useNavigate } from "react-router-dom";
import * as yup from 'yup' 
import { Formik, Form, Field, ErrorMessage } from 'formik';

const RegisterForm = () => {

  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const passwordValidator =  yup.string().min(8, "password is too short")
  .matches( /[a-z]+/, "needs to contain lowercase letter")
  .matches( /[A-Z]+/, "needs to contain uppercase letter")
  .matches( /[0-9]+/, "needs to contain number letter")
  .required();

  const phoneValidator = yup.number()
  .min(Math.pow(10, 6), "Must have at least 6 digits")
  .max(Math.pow(10, 12), "Must have less than 12 digits")
  .typeError("Must be a phone number").required();

  const schema = yup.object().shape({
    "first name": yup.string().required(),
    "last name": yup.string().required(),
    "phone number": phoneValidator,
    email: yup.string().email().required(),
    password: passwordValidator,
    "confirm password": passwordValidator.oneOf([yup.ref('password')], 'Passwords must match'),
  })

  const navigate = useNavigate();

  const onClick = () => {
    (async function () {
      try {
        await AuthService.register({ Email: email, Password: password, Phone: phoneNumber, FirstName: firstName, LastName: lastName })
        navigate("/login")
      } catch (error: any) {
       // alert(error.response.data);
      }
    })()

  }

  return (
    <Formik
    initialValues={{
      "first name": "",
      "last name": "",
      "phone number": "",
      email: "",
      password: "",
      "confirm password": "",
    }}
    validationSchema={schema}
    onSubmit={values => {

    }}
  >
    {({ errors, touched, setFieldValue }) => (
      <Form>
        <Field name="first name" component={ InputField} className={RegisterFormCSS.inlineInput} usage="First name" value={firstName} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                    setFirstName(e.target.value);
                    setFieldValue("first name", e.target.value);
                  }}/>
        <ErrorMessage name="first name" />

        <Field component={ InputField} className={`alignRight ${RegisterFormCSS.inlineInput}`}  usage="Last name" value={lastName} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                    setLastName(e.target.value);
                    setFieldValue("last name", e.target.value);
                  }}/>
        <ErrorMessage name="last name" />

        <Field component={ InputField} className={RegisterFormCSS.inlineInput} usage="Email" value={email} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                    setEmail(e.target.value);
                    setFieldValue("email", e.target.value);
                  }}/>
        <ErrorMessage name="email" />

        <Field component={ InputField} className={`alignRight ${RegisterFormCSS.inlineInput}`} usage="Phone number" value={phoneNumber} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                    setPhoneNumber(e.target.value);
                    setFieldValue("phone number", e.target.value);
                  }}/>
        <ErrorMessage name="phone number" />

        <Field component={ InputField} className={RegisterFormCSS.inlineInput} usage="Password" value={password} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                    setPassword(e.target.value);
                    setFieldValue("password", e.target.value);
                  }}/>
        <ErrorMessage name="password" />

        <Field component={ InputField} className={`alignRight ${RegisterFormCSS.inlineInput}`} usage="Confirm password" value={confirmPassword} onChange={(e:React.ChangeEvent<HTMLInputElement>) => {
                    setConfirmPassword(e.target.value);
                    setFieldValue("confirm password", e.target.value);
                  }}/>
        <ErrorMessage name="confirm password" />


        <div className={RegisterFormCSS.button}>
          <span className="alignRight">
            <Button submit="submit" onClick={onClick} text="Get started" />
          </span>
        </div>
      </Form>
  )}
  </Formik >
)}
export default RegisterForm