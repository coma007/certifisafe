import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import RegisterFormCSS from './RegisterForm.module.scss'
import { useRef, useState } from 'react';
import { AuthService } from 'features/auth/services/AuthService';
import { useNavigate } from "react-router-dom";
import * as yup from 'yup'
import { Formik, Form, Field, ErrorMessage } from 'formik';
import ReCAPTCHA from 'react-google-recaptcha';
import ErrorMsg from 'components/error/ErrorMsg';

const RegisterForm = () => {

  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const passwordValidator = yup.string().min(8, "password is too short")
    .matches(/[a-z]+/, "password needs to contain lowercase letter")
    .matches(/[A-Z]+/, "password needs to contain uppercase letter")
    .matches(/[0-9]+/, "password needs to contain number letter")
    .required();

  const phoneValidator = yup.number()
    .min(Math.pow(10, 6), "phone needs to have at least 6 digits")
    .max(Math.pow(10, 12), "phone needs to have less than 12 digits")
    .typeError("must be a phone number").required();

  const schema = yup.object().shape({
    "first name": yup.string().required(),
    "last name": yup.string().required(),
    "phone number": phoneValidator,
    email: yup.string().email().required(),
    password: passwordValidator,
    "confirm password": passwordValidator.oneOf([yup.ref('password')], 'Passwords must match'),
  })

  const navigate = useNavigate();
  const captchaRef: any = useRef(null)

  const handleSubmit = async () => {
    const token = captchaRef.current?.getValue();
    captchaRef.current.reset();
    try {
      await AuthService.register({ Email: email, Password: password, Phone: phoneNumber, FirstName: firstName, LastName: lastName, Token: token })
      navigate("/login")
    } catch (error: any) {
      alert(error.response.data);
    }
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
        token: "",
      }}
      validateOnChange
      validationSchema={schema}
      onSubmit={handleSubmit}
    >
      {({ errors, touched, setFieldValue, validateForm, isValid }) => (
        <Form>
          <div className={RegisterFormCSS.grid}>
            <div>
              <Field name="first name" component={InputField} className={RegisterFormCSS.inlineInput} usage="First name" value={firstName} onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                setFirstName(e.target.value);
                setFieldValue("first name", e.target.value);
              }} />
              <ErrorMsg val={errors["first name"]} />
            </div>
            <div>
              <Field component={InputField} className={`alignRight ${RegisterFormCSS.inlineInput}`} usage="Last name" value={lastName} onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                setLastName(e.target.value);
                setFieldValue("last name", e.target.value);
              }} />
              <ErrorMsg val={errors["last name"]} customClass="alignRight" />
            </div>
            <div>
              <Field component={InputField} className={RegisterFormCSS.inlineInput} usage="Email" value={email} onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                setEmail(e.target.value);
                setFieldValue("email", e.target.value);
              }} />
              <ErrorMsg val={errors["email"]} />
            </div>
            <div>
              <Field component={InputField} className={`alignRight ${RegisterFormCSS.inlineInput}`} usage="Phone number" value={phoneNumber} onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                setPhoneNumber(e.target.value);
                setFieldValue("phone number", e.target.value);
              }} />
              <ErrorMsg val={errors["phone number"]} customClass="alignRight" />
            </div>
            <div>
              <Field component={InputField} className={RegisterFormCSS.inlineInput} usage="Password" value={password} onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                setPassword(e.target.value);
                setFieldValue("password", e.target.value);
              }} />
              <ErrorMsg val={errors["password"]} />
            </div>
            <div>
              <Field component={InputField} className={`alignRight ${RegisterFormCSS.inlineInput}`} usage="Confirm password" value={confirmPassword} onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                setConfirmPassword(e.target.value);
                setFieldValue("confirm password", e.target.value);
              }} />
              <ErrorMsg val={errors["confirm password"]} customClass="alignRight" />
            </div>
          </div>


          <div className={RegisterFormCSS.button}>
            <span className="alignRight">
              <Button submit="submit" onClick={null} text="Get started" />
            </span>
          </div>
          <ReCAPTCHA className='recaptcha' sitekey={process.env.REACT_APP_SITE_KEY as string} ref={captchaRef} />
        </Form>
      )}
    </Formik >
  )
}
export default RegisterForm