import LoginForm from 'features/auth/components/LoginForm/LoginForm'
import Banner from 'components/navigation/Banner/Banner'
import { AuthService } from 'features/auth/services/AuthService'

import TooltipCSS from 'components/view/Tooltip/Tooltip.module.scss'

import Gmail from 'assets/oauth/gmail.png'
import Facebook from 'assets/oauth/facebook.png'
import Tooltip from 'components/view/Tooltip/Tooltip'
import { useState } from 'react'
import TwoFactorForm from 'features/auth/components/TwoFactorForm/TwoFactorForm'

const LoginPage = () => {
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
                        <h2>Sign in</h2>
                        <span >
                            Welcome back!
                            <br />Please enter your login details to access your account.
                        </span>
                    </div>
                    <LoginForm twoFactor={sendCode} />
                    <div className="oauth">
                        Or use alternative way to sign in <br />
                        <button className={TooltipCSS.bottomTooltip}>
                            <img src={Gmail} />
                            <Tooltip tooltipText="Sign in with Gmail account" />
                        </button>
                        <button className={TooltipCSS.bottomTooltip}>
                            <img src={Facebook} />
                            <Tooltip tooltipText="Sign in with Facebook account" />
                        </button>
                    </div>
                    <div className="authBottomMessage">
                        Do not have an account ?
                        <br /> <a href='register'>Sign up here.</a>
                    </div>
                </div >
            ) : (
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
            )}
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
