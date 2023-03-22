import React from 'react'
import MenuItemCSS from './MenuItem.module.scss'

const MenuItem = ({ image, className }: { image: string, className: string }) => {
    return (
        <a >
            <img className={`${MenuItemCSS.item} ${className}`} src={image} />
        </a>
    )
}

export default MenuItem