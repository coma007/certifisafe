import Banner from 'components/navigation/Banner/Banner';
import EmailForm from 'features/auth/components/EmailForm/EmailForm';
import PasswordRecoveryForm from 'features/auth/components/PasswordRecoveryForm/PasswordRecoveryForm';
import { AuthService } from 'features/auth/services/AuthService';
import { send } from 'process';
import React, { useState } from 'react'

const PasswordRecoveryPage = () => {
  let [isCodeSent, setIsCodeSent] = useState<boolean>(false);

  const sendCode = () => {
    setIsCodeSent(true);
  }

  return (
    <div className="page pageTwoCols">
      <div>
        <Banner />
      </div>
      {isCodeSent === false ? (
        <div className="rightCol">
          <div className="authTitle">
            <h2>Forgot password ?</h2>
            <span>
              Please enter your email, and shortly, we will send you verification code for password reset.
            </span>
          </div>
          <EmailForm onClick={sendCode} />
        </div>
      ) : (
        <div className="rightCol">
          <div className="authTitle">
            <h2>Reset password</h2>
            <span>
              We have sent you email and SMS with a verification code.
              <br />
              Please enter the code below and a new password.
            </span>
          </div>
          <PasswordRecoveryForm resetPage={null}/>
        </div>
      )}
    </div>
  )
}

export default PasswordRecoveryPage