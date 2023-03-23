import CardCSS from "./Card.module.scss"

const Card = ({ children }: { children: React.ReactNode }) => {
    return (
        <a className={CardCSS.card}>{children}</a>
    )
}

export default Card