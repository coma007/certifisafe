import { User } from "features/auth/types/User"

export interface Request {
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