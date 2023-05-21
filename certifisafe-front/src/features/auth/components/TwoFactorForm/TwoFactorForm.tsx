import Button from 'components/forms/Button/Button'
import React from 'react'
import VerificationInput from 'react-verification-input'
import TwoFactorFormCSS from './TwoFactorForm.module.scss'

const TwoFactorForm = () => {
    return (
        <div>
            <VerificationInput length={4} placeholder={""} autoFocus={true} classNames={{
                container: TwoFactorFormCSS.container,
                character: "codeField",
                characterInactive: "codeFieldInactive",
                characterSelected: "codeFieldActive",
            }} />
            <span className="alignRight">
                <Button onClick={undefined} text="Sign in" />
            </span>
        </div>
    )
}

export default TwoFactorForm