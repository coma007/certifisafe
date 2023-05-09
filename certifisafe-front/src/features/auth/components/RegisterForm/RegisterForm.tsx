import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import RegisterFormCSS from './RegisterForm.module.scss'
import { useState } from 'react';
import { AuthService } from 'features/auth/services/AuthService';

const RegisterForm = () => {
  
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const firstNameChange = (event : any) => {
    setFirstName(event.target.value);
  };

  const lastNameChange = (event : any) => {
    setLastName(event.target.value);
  };

  const phoneNumberChange = (event : any) => {
    setPhoneNumber(event.target.value);
  };

  const emailChange = (event : any) => {
    setEmail(event.target.value);
  };

  const passwordChange = (event : any) => {
    setPassword(event.target.value);
  };

  const onClick = () => {
    (async function () {
      try {
          await AuthService.register({Email: email, Password: password, Phone: phoneNumber, FirstName: firstName, LastName: lastName})
          alert("Success")
      } catch (error: any) {
          alert(error.response.data);
      }
  })()
    
  }

  return (
    <div className={RegisterFormCSS.form}>
      <InputField className={RegisterFormCSS.inlineInput} usage="First name" value={firstName} onChange={firstNameChange}/>
      <InputField className={`alignRight ${RegisterFormCSS.inlineInput}`} usage="Last name" value={lastName} onChange={lastNameChange}/>
      <InputField className={RegisterFormCSS.input} usage="Phone number" value={phoneNumber} onChange={phoneNumberChange}/>
      <InputField className={RegisterFormCSS.input} usage="Email" value={email} onChange={emailChange}/>
      <InputField className={RegisterFormCSS.input} usage="Password" value={password} onChange={passwordChange}/>
      <div className={RegisterFormCSS.button}>
        <span className="alignRight">
          <Button onClick={onClick} text="Get started" />
        </span>
      </div>
    </div >
  )
}

export default RegisterForm