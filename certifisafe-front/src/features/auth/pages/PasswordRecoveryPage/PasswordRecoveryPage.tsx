import Button from 'components/forms/Button/Button';
import InputField from 'components/forms/InputField/InputField';
import Banner from 'components/navigation/Banner/Banner';
import PasswordRecoveryForm from 'features/auth/components/PasswordRecoveryForm/PasswordRecoveryForm';
import React, { useState } from 'react'
import VerificationInput from "react-verification-input";

const PasswordRecoveryPage = () => {

  return (

    <div className="page pageTwoCols">
      <div>
        <Banner />
      </div>
      <div className="rightCol">
        <div className="authTitle">
          <h2>Reset password</h2>
          <span >
            We have sent you an email with verification code.
            <br />
            Please enter the code bellow and new password.
          </span>
        </div>
        <PasswordRecoveryForm />
      </div >
    </div>
  )
}

export default PasswordRecoveryPage