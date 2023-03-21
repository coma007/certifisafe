import React from 'react'
import BannerCSS from "./Banner.module.scss"
import Logo from "assets/logo-full-white.png"

const Banner = () => {
    return (
        <div className={`half-screen-component ${BannerCSS.banner}`} >
        <img className={BannerCSS.logo} src={Logo}></img>
        </div >
    )
}

export default Banner