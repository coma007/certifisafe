import Banner from 'components/navigation/Banner/Banner'
import RegisterForm from 'features/auth/components/RegisterForm/RegisterForm'
import TooltipCSS from 'components/view/Tooltip/Tooltip.module.scss'

import Gmail from 'assets/oauth/gmail.png'
import Facebook from 'assets/oauth/facebook.png'
import Tooltip from 'components/view/Tooltip/Tooltip'

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
        <div className="oauth">
          Or use alternative way to sign up <br />
          <button className={TooltipCSS.bottomTooltip}>
            <img src={Gmail} />
            <Tooltip tooltipText="Sign up with Gmail account" />
          </button>
          <button className={TooltipCSS.bottomTooltip}>
            <img src={Facebook} />
            <Tooltip tooltipText="Sign up with Facebook account" />
          </button>
        </div>
        <div className="authBottomMessage">
          Already have an account ?
          <br /> <a href='login'>Sign in here.</a>
        </div>
      </div >
    </div>
  )
}

export default RegisterPage