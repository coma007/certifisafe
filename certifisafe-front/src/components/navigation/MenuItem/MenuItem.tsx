import MenuItemCSS from './MenuItem.module.scss'
import Tooltip from '../../view/Tooltip/Tooltip'
import TooltipCSS from '../../view/Tooltip/Tooltip.module.scss'
import { Link } from "react-router-dom";

interface IMenuItemProps { image: string, className: string, tooltipText: string, tooltip: boolean, path: string }

const MenuItem = ({ image, className, tooltipText, tooltip, path }: IMenuItemProps) => {
    return (
        <Link to={path}
            className={`${tooltip ? TooltipCSS.tooltip : TooltipCSS.noTooltip}`}>
            <img className={`${MenuItemCSS.item} ${className}`} src={image} />
            <Tooltip tooltipText={tooltipText} />
        </Link>
    )
}

export default MenuItem