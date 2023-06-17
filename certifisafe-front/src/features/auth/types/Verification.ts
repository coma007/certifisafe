export interface Code {
    VerificationCode: string
}

export interface PasswordResetRequest {
    Email: string,
    Type: number
}

export interface PasswordReset {
    VerificationCode: string,
    newPassword: string
}