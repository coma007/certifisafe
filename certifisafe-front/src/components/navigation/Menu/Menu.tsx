import MenuItem from '../MenuItem/MenuItem'
import MenuCSS from './Menu.module.scss'

import Logo from "assets/menu/logo.png"
import Home from "assets/menu/home-white.png"
import Certificate from "assets/menu/certificate-white.png"
import Verify from "assets/menu/verify-white.png"
import Create from "assets/menu/create-white.png"
import Request from "assets/menu/request-white.png"
import Profile from "assets/menu/profile-white.png"
import Logout from "assets/menu/logout-white.png"

const Menu = () => {
    return (
        <div className={MenuCSS.menu}>
            <MenuItem className={MenuCSS.logo} tooltipText="" tooltip={false} image={Logo} path="/" />
            <MenuItem className={MenuCSS.nonMainOption} tooltipText="Homepage" tooltip={true} image={Home} path="/" />
            <hr className={MenuCSS.separator} />
            <MenuItem className={MenuCSS.margin} tooltipText="Certificates overview" tooltip={true} image={Certificate} path="/certificates" />
            <MenuItem className={MenuCSS.margin} tooltipText="Verify certificate" tooltip={true} image={Verify} path="/verify" />
            <MenuItem className={MenuCSS.margin} tooltipText="New certificate" tooltip={true} image={Create} path="" />
            <MenuItem className={MenuCSS.margin} tooltipText="Requests" tooltip={true} image={Request} path="/requests" />
            <MenuItem className={MenuCSS.nonMainOption} tooltipText="Profile" tooltip={true} image={Profile} path="" />
            <MenuItem className={MenuCSS.nonMainOption} tooltipText="Sign out" tooltip={true} image={Logout} path="/logout" />
        </div>
    )
}

export default Menu