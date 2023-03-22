import React from 'react'
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
            <MenuItem className={MenuCSS.logo} image={Logo}></MenuItem>
            <MenuItem className={MenuCSS.nonMainOption} image={Home}></MenuItem>
            <hr className={MenuCSS.separator} />
            <MenuItem className={MenuCSS.margin} image={Certificate}></MenuItem>
            <MenuItem className={MenuCSS.margin} image={Verify}></MenuItem>
            <MenuItem className={MenuCSS.margin} image={Create}></MenuItem>
            <MenuItem className={MenuCSS.margin} image={Request}></MenuItem>
            <MenuItem className={MenuCSS.nonMainOption} image={Profile}></MenuItem>
            <MenuItem className={MenuCSS.nonMainOption} image={Logout}></MenuItem>
        </div>
    )
}

export default Menu