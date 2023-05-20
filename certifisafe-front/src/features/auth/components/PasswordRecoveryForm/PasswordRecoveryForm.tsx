import Button from 'components/forms/Button/Button'
import InputField from 'components/forms/InputField/InputField'
import React from 'react'
import VerificationInput from 'react-verification-input'
import PasswordRecoveryFormCSS from "./PasswordRecoveryForm.module.scss"


const PasswordRecoveryForm = () => {
    return (
        <div>
            <VerificationInput length={4} placeholder={""} autoFocus={true} classNames={{
                character: PasswordRecoveryFormCSS.field,
                characterInactive: PasswordRecoveryFormCSS.fieldInactive,
                characterSelected: PasswordRecoveryFormCSS.fieldActive,
            }} />
            <p className={PasswordRecoveryFormCSS.label}>
                Did not receive code ?
                <a onClick={undefined}>
                    <b> Click to resend.</b>
                </a>
            </p>
            <InputField usage='New password' className={PasswordRecoveryFormCSS.input} />
            <InputField usage='Confirm password' className={PasswordRecoveryFormCSS.input} />
            <span className="alignRight">
                <Button onClick={undefined} text="Reset" />
            </span>
        </div>
    )
}

export default PasswordRecoveryForm