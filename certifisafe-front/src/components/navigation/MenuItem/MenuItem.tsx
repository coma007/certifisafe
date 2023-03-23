import React from 'react'
import MenuItemCSS from './MenuItem.module.scss'

const MenuItem = ({ image, className, tooltipText, tooltip }:
    { image: string, className: string, tooltipText: string, tooltip: boolean }) => {
    return (
        <a className={`${tooltip ? "tooltip" : "noTooltip"}`}>
            <img className={`${MenuItemCSS.item} ${className}`} src={image} />
            <span className="tooltiptext">{tooltipText}</span>
        </a>
    )
}

export default MenuItem