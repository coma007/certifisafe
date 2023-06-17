import LoginForm from 'features/auth/components/LoginForm/LoginForm'
import Banner from 'components/navigation/Banner/Banner'
import { AuthService } from 'features/auth/services/AuthService'

import TooltipCSS from 'components/view/Tooltip/Tooltip.module.scss'

import Gmail from 'assets/oauth/gmail.png'
import Facebook from 'assets/oauth/facebook.png'
import Tooltip from 'components/view/Tooltip/Tooltip'
import { useState } from 'react'
import TwoFactorForm from 'features/auth/components/TwoFactorForm/TwoFactorForm'
import { useLocation } from 'react-router-dom'
import PasswordRecoveryForm from 'features/auth/components/PasswordRecoveryForm/PasswordRecoveryForm'

const LoginPage = () => {
    let [isCodeSent, setIsCodeSent] = useState<boolean>(false);
    let [isPasswordReset, setIsPasswordReset] = useState<boolean>(false);
    let [isBasePage, setIsBasePage] = useState<boolean>(true);

    const sendCode = () => {
        setIsCodeSent(true);
        setIsBasePage(false);
    }
    
    const resetPassword = () => {
        setIsPasswordReset(true);
        setIsBasePage(false);
    }

    const resetPage = () => {
        setIsBasePage(true);
        setIsPasswordReset(false);
    }

    const oauth = () => {
        (async function () {
            try {
                await AuthService.oauth();
            } catch (error: any) {
                alert(error.response.data);
            }
        })()
    }

    return (
        <div className="page pageTwoCols">
            <div>
                <Banner />
            </div>
            {isBasePage === true ? (
                <div className="rightCol">
                    <div className="authTitle">
                        <h2>Sign in</h2>
                        <span >
                            Welcome back!
                            <br />Please enter your login details to access your account.
                        </span>
                    </div>
                    <LoginForm twoFactor={sendCode} resetPassword={resetPassword} />
                    <div className="oauth" onClick={oauth}>
                        Or use alternative way to sign in <br />
                        <button className={TooltipCSS.bottomTooltip} >
                            <img src={Gmail} />
                            <Tooltip tooltipText="Sign in with Gmail account" />
                        </button>
                    </div>
                    <div className="authBottomMessage">
                        Do not have an account ?
                        <br /> <a href='register'>Sign up here.</a>
                    </div>
                </div >
            ) : null}
            { isCodeSent === true ? (
                <div className="rightCol">
                    <div className="authTitle">
                        <h2>Confirm it is you</h2>
                        <span >
                            We have sent you email and SMS with a verification code.
                            <br />
                            Please enter the code below to confirm your identity.
                        </span>
                    </div>
                    <TwoFactorForm />
                </div >
            ) : null}
            { isPasswordReset === true ? (
            <div className="rightCol">
                <div className="authTitle">
                    <h2>Renew password</h2>
                    <span>
                        Your password is too old, you need to renew it now.
                        <br />
                        We have sent you email and SMS with a verification code.
                        <br />
                        Please enter the code below and a new password.
                    </span>
                </div>
                <PasswordRecoveryForm resetPage={resetPage} />
            </div>
            ) : null}
        </div>
    )
}

export default LoginPage

export const Logout = () => {

    AuthService.logout()
    window.location.href = "/login";

    return (
        <></>
    )
}
