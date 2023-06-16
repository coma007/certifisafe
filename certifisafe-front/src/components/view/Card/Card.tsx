import { Link } from "react-router-dom"
import CardCSS from "./Card.module.scss"

const Card = (props: { children: React.ReactNode, link?: string, onClick? : any }) => {
    return (
        <a href={props.link} className={CardCSS.card} onClick={props.onClick} > {props.children}</a >
    )
}

export default Card