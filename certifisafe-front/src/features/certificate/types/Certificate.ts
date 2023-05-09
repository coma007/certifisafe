import { User } from "features/auth/types/User"

export interface Certificate {
    Serial: number
    Name: string
    ValidFrom: Date
    ValidTo: Date
    Issuer: User
    Subject: User
    Status: string
    Type: string
}