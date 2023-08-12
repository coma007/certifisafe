# CertifiSafe
App for certificates managing.


### Team members
[Nemanja Majstorović](github.com/Nemanja3314)  
[Nemanja Dutina](github.com/eXtremeNemanja)  
[Milica Sladaković](github.com/coma007)

## Table of Contents
- [Overview](#overview)
- [Understanding x509 Certificates](#understanding-x509-certificates)
  - [What are x509 Certificates ?](#what-are-x509-certificates)
  - [x509 Certificate Hierarchy](#x509-certificate-hierarchy)
  - [Components of an x509 Certificate](#components-of-an-x509-certificate)
- [Certificate Management](#certificate-management)
- [User Types](#user-types)
- [Logging and Security](#logging-and-security)
- [Integrations](#integrations)
- [Installation and Usage](#installation-and-usage)

## Overview

CertifiSafe is a certificate management web application designed to enable the creation, management, and validation of certificates. Using x509 technology, CertifiSafe simplifies certificate handling by utilizing a monolithic architecture, backed by PostgreSQL for data storage, Golang for secure backend operations, and React for frontend user interactions.


## Understanding x509 Certificates

In the world of digital security, x509 certificates serve as a fundamental building block, providing a standardized format for ensuring the authenticity and security of digital communications. These certificates play main role in this application by guaranteeing the integrity of the certificate ecosystem.

### What are x509 Certificates?

x509 certificates, also known as SSL/TLS certificates or public key certificates, are a digital form of identification used to authenticate and encrypt communications over the internet. They are primarily utilized to establish secure connections between servers and clients, safeguarding sensitive data from eavesdropping.

### x509 Certificate Hierarchy

Certificates are organized into a hierarchical structure to establish trust and enable secure communication. There are three distinct types of certificates: root, intermediate, and end certificates.

- **Root Certificates** - Root certificates are the foundational certificates at the top of the certificate hierarchy. These certificates are self-signed and serve as the ultimate trust anchors for the entire system. Root certificates are typically provided by trusted Certificate Authorities (CAs) and establish the trustworthiness of the entire certificate chain.
- **Intermediate Certificates** - Intermediate certificates sit between root and end certificates, forming a bridge of trust. They are issued by the root certificates and provide an additional layer of security by preventing direct access to root certificates. 
- **End Certificates** - End certificates, also known as leaf certificates, are the certificates issued to end entities, such as websites or users. These certificates contain the public key of the entity, along with identifying information. End certificates are signed by intermediate certificates, forming a chain of trust back to the root certificate. They are used to secure communication and validate the authenticity of entities within the system.

In this project, the management of these three types of certificates is central to the platform's functionality, ensuring the security and integrity of all types of interactions.

### Components of an x509 Certificate

An x509 certificate comprises several essential components:

- **Subject**: Identifies the entity the certificate is issued to, often in the form of a domain name.
- **Issuer**: Specifies the entity that issued the certificate, typically a Certificate Authority (CA).
- **Public Key**: The encryption key used for secure communication.
- **Digital Signature**: A cryptographic signature created by the issuer to validate the certificate's authenticity.
- **Validity Period**: Indicates the duration during which the certificate is considered valid.
- **Key Usage**: Defines the allowed operations that can be performed using the public key.
- **Extended Key Usage**: Specifies the intended uses for the certificate.
- **Thumbprint**: A unique hash value that identifies the certificate.
- **Version**: Denotes the x509 certificate version.


## Certificate Management

CertifiSafe ensures certificates are structured in a secure and organized manner:

- **x509 Certificates**: The application uses the x509 standard for certificates, ensuring compatibility and adherence to best practices.
- **Public and Private Storage**: Certificates are stored in the `public` and `private` storages. The `public` directory contains the `.crt` public part of the certificate, while the `private` directory contains the `.key` private part of the certificate.
- **Database Metadata**: To streamline access and retrieval, certificate metadata is stored in the PostgreSQL database.
- **Certificate Operations**: Users can perform a range of operations, including listing, downloading, validating, creating, withdrawing, and managing requests for root, intermediate, and end certificates.

## User Types

- **Unauthenticated User**: These users can easily register or log in if already registered, gaining access to their personal certificate management dashboard.
- **Authenticated User**: Authenticated users can explore, download, validate, and request creation of certificates.
- **Admin**: Admin accounts have elevated privileges, including the management of root certificates and the approval of certificate issuance requests.

## Logging and Security

- **Data Validation**: Data validation, including defense against injection, XSS, and path traversal attacks is implemented. React ensures security by promoting data and rendering separation to resist injection and XSS attacks, while Golang with GORM fights path traversal vulnerabilities through parameterized queries, ORM abstraction, prepared statements, and input validation, enabling secure database interactions. File uploads are also restricted by size.
- **Two-Factor Authentication (2FA)**: App enforces 2FA during login, enhancing security by requiring users to authenticate via a verification code sent to their selected device.
- **Password Rotation**: Users are prompted to change their passwords after a defined time, enhancing security, and preventing password reuse.
- **Custom Logging Middleware**: A custom middleware implemented in Golang provides detailed logging capabilities, ensuring accountability and aiding in debugging.


## Integrations

- **Google Sign-In/Sign-Up**: OAuth-based authentication streamlines user registration and login, reducing friction and enhancing security.
- **Mailgun Integration**: CertifiSafe utilizes Mailgun as the mail server for efficient email communication.
- **Twilio Integration**: Twilio integration facilitates the sending of messages for authorization purposes.


## Installation and Usage

1. **Clone the repository**: Start by cloning the CertifiSafe repository: `git clone https://github.com/coma007/certifisafe.git`
2. **Navigate to the project directory**: Move into the cloned directory: `cd certifisafe`
3. **Database setup**: Initialize the database schema by executing the provided [SQL scripts](https://github.com/coma007/certifisafe/tree/master/certifisafe-back/resources/database).
4. **Backend configuration**: Navigate to the `certifisafe-back` directory and create configuration file `config.yaml` with the same structure as shown in [`config-structure.yaml`](https://github.com/coma007/certifisafe/blob/master/certifisafe-back/config-structure.yaml).
5. **Generate root certificate**: Call function for generating the root certificate in [`main.go`](https://github.com/coma007/certifisafe/blob/master/certifisafe-back/main.go). The root certificate will be stored in `public/server.crt` and `private/server.key`.
6. **Start the backend**: Run the backend server using `go run main.go`.
7. **Frontend setup**: Navigate to the `certifisafe-front` directory and install dependencies using `npm install`. Setup the configuration in `.env` file (contains `RECAPTCHA_SECRET_KEY` and `RECAPTCHA_SITE_KEY` for ReCaptha config).
8. **Add the server certificate to frontend**: Copy the `certificate-back/public/server.crt` into `certificate-front/cert/cert.crt`. Copy the `certificate-back/private/server.key` into `certificate-front/cert/key.key`. Those are server generated certificates that will enable secure HTTPS communication.
9. **Start the frontend**: Start the frontend using `npm start`.


## Happy Certificate Managing ! 🛡️
