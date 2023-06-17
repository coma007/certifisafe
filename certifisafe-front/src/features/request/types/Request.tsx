import { User } from "features/auth/types/User"

export interface Request {
    ID: number
    CertificateName: string
    Date: Date
    Subject: User
    Status: string
    CertificateType: string
}

export interface CreateRequestDTO {
    CertificateName: string
    CertificateType: string
    Token: string
    ParentSerial: number
}