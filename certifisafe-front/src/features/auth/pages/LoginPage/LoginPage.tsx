import LoginForm from 'features/auth/components/LoginForm/LoginForm'
import Banner from 'components/navigation/Banner/Banner'
import { AuthService } from 'features/auth/services/AuthService'
import { Navigate } from 'react-router-dom'

const LoginPage = () => {
    return (
        <div className="page pageTwoCols">
            <div>
                <Banner />
            </div>
            <div className="rightCol">
                <div className="authTitle">
                    <h2>Sign in</h2>
                    <span >
                        Welcome back!
                        <br />Please enter your login details to access your account.
                    </span>
                </div>
                <LoginForm />
                <div className="authBottomMessage">
                    Do not have an account ?
                    <br /> <a href='register'>Sign up here.</a>
                </div>
            </div >
        </div>
    )
}

export default LoginPage

export const Logout = () => {

    AuthService.logout()
    window.location.href = "/";

    return (
        <></>
    )
}
