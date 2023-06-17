import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';

function RedirectPage() {
    useEffect(() => {
        const searchParams = new URLSearchParams(window.location.search);
        const token = searchParams.get('token');
        if (token) {
            localStorage.setItem('token', token);
        }
        window.location.href = "/";
    }, []);

    return (
        <div>
            <p>Redirecting...</p>
        </div>
    );
}

export default RedirectPage;

