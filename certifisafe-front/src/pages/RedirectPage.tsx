import { useEffect } from 'react';

function RedirectPage() {
    useEffect(() => {
        const token = getCookieValue('token');

        if (typeof window !== 'undefined') {
            localStorage.setItem('token', token);
        }

        window.location.href = "/";
    }, []);

    const getCookieValue = (name: any) => {
        if (typeof document !== 'undefined') {
            const cookies = document.cookie.split('; ');
            for (const element of cookies) {
                const cookie = element.split('=');
                if (cookie[0] === name) {
                    return cookie[1];
                }
            }
        }
        return '';
    };

    return (
        <div>
            <p>Redirecting...</p>
        </div>
    );
}

export default RedirectPage;

