export interface User {
    Email: string
    FirstName: string
    LastName: string
    Phone: string
}

export interface Credentials {
    Email: string
    Password: string
    Token?: string
}

export interface UserRegister {
    Email: string
    Password: string
    FirstName: string
    LastName: string
    Phone: string
    Token?: string
}
