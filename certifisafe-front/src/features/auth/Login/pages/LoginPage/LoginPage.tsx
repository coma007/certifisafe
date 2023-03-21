import React from 'react'
import LoginForm from '../../components/LoginForm/LoginForm'
import LoginPageCSS from './LoginPage.module.scss'
import Banner from 'components/navigation/Banner/Banner'

const LoginPage = () => {
    return (
        <div className={`page-two-cols ${LoginPageCSS.page}`}>
            <div>
                <Banner></Banner>
            </div>
            <div className={LoginPageCSS.rightCol}>
                <div className={LoginPageCSS.title}>
                    <h2>Sign in</h2>
                    <span >Welcome back!<br />Please enter your login details to access your account.</span>
                </div>
                <LoginForm />
                <div className={LoginPageCSS.bottomMessage}>
                    Do not have an account ?<br /> <a href='#'>Sign up here.</a>
                </div>
            </div >
        </div>
    )
}

export default LoginPage