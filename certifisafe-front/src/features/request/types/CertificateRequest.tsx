import { User } from "features/auth/types/User"

export interface CertificateRequest {
    CertificateName: string
    Date: Date
    Subject: User
    Status: string
    CertificateType: string
}