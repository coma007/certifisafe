import Button from 'components/forms/Button/Button'
import React, { useState } from 'react'
import VerificationInput from 'react-verification-input'
import TwoFactorFormCSS from './TwoFactorForm.module.scss'
import { AuthService } from 'features/auth/services/AuthService'
import { redirect, useNavigate } from 'react-router-dom'


const TwoFactorForm = () => {
    const navigate = useNavigate();
    const [code, setCode] = useState('');

    const onClick = () => {
        (async function () {
            try {
                const jwt = await AuthService.tfactorauth({ VerificationCode: code });
                localStorage.setItem("token", jwt)
                navigate("/")
            }catch (error: any) {
                alert(error.response.data);
            }
        })()
    }


    return (
        <div>
            <VerificationInput value={code} onChange={(e:string) => {
                    setCode(e);
                  }} length={4} placeholder={""} autoFocus={true} classNames={{
                container: TwoFactorFormCSS.container,
                character: "codeField",
                characterInactive: "codeFieldInactive",
                characterSelected: "codeFieldActive",
            }} />
            <span className="alignRight">
                <Button submit={undefined} onClick={onClick} text="Sign in" />
            </span>
        </div>
    )
}

export default TwoFactorForm