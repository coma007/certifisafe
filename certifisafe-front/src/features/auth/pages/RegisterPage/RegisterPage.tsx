import Banner from 'components/navigation/Banner/Banner'
import RegisterForm from 'features/auth/components/RegisterForm/RegisterForm'

const RegisterPage = () => {
  return (
    <div className="page pageTwoCols">
      <div>
        <Banner />
      </div>
      <div className="rightCol">
        <div className="authTitle">
          <h2>Sign up</h2>
          <span >
            Welcome !
            <br />Sign up now to get started.
          </span>
        </div>
        <RegisterForm />
        <div className="authBottomMessage">
          Do not have an account ?
          <br /> <a href='#'>Sign up here.</a>
        </div>
      </div >
    </div>
  )
}

export default RegisterPage