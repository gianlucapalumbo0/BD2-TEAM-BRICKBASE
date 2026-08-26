console.log("VERSIONE NUOVA CARICATA");

document.addEventListener('DOMContentLoaded', () => {
    const btn = document.getElementById('mobile-menu-btn');
    const menu = document.getElementById('mobile-menu');
    const iconOpen = document.getElementById('icon-open');
    const iconClose = document.getElementById('icon-close');

    if (btn) {
        btn.addEventListener('click', () => {
            menu.classList.toggle('hidden');
            iconOpen.classList.toggle('hidden');
            iconClose.classList.toggle('hidden');
        });
    }

    updateNavbarAuthState();
    setupLogoutHandlers();
});

function formatRole(role) {
    const labels = {
        ADMIN: 'Admin',
        USER: 'Utente'
    };
    return labels[role] || role;
}

async function fetchCurrentUser() {
    try {
        const response = await fetch('/api/me', {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            const user = await response.json();
            localStorage.setItem('user', JSON.stringify(user));
            return user;
        }

        localStorage.removeItem('user');
        return null;
    } catch {
        const raw = localStorage.getItem('user');
        if (!raw) return null;

        try {
            return JSON.parse(raw);
        } catch {
            localStorage.removeItem('user');
            return null;
        }
    }
}

function renderNavbarAuthState(user) {
    const loginBtn = document.getElementById('nav-login-btn');
    const registerBtn = document.getElementById('nav-register-btn');
    const userInfo = document.getElementById('nav-user-info');
    const userName = document.getElementById('nav-user-name');
    const userRole = document.getElementById('nav-user-role');
    const logoutBtn = document.getElementById('nav-logout-btn');

    const mobileGuest = document.getElementById('mobile-auth-guest');
    const mobileUser = document.getElementById('mobile-auth-user');
    const mobileUserName = document.getElementById('mobile-user-name');
    const mobileUserRole = document.getElementById('mobile-user-role');

    // --- GESTIONE DINAMICA PULSANTE "INSERISCI" (ADMIN) ---
    // Rimuoviamo eventuali link admin esistenti per evitare duplicati
    document.querySelectorAll('.dynamic-admin-link').forEach(el => el.remove());

    if (user) {
        const fullName = `${user.first_name || ''} ${user.last_name || ''}`.trim();
        const roleLabel = formatRole(user.role);

        loginBtn?.classList.add('hidden');
        registerBtn?.classList.add('hidden');
        userInfo?.classList.remove('hidden');
        logoutBtn?.classList.remove('hidden');

        if (userName) userName.textContent = fullName;
        if (userRole) userRole.textContent = roleLabel;

        mobileGuest?.classList.add('hidden');
        mobileUser?.classList.remove('hidden');
        if (mobileUserName) mobileUserName.textContent = fullName;
        if (mobileUserRole) mobileUserRole.textContent = roleLabel;

        // Se l'utente è un ADMIN, iniettiamo il link "Inserisci"
        if (user.role && user.role.toUpperCase() === 'ADMIN') {
            // 1. Inserimento per Desktop (cerchiamo il blocco info utente o i pulsanti di navigazione)
            if (userInfo && userInfo.parentNode) {
                const adminDesktopLink = document.createElement('a');
                adminDesktopLink.href = '/admin';
                adminDesktopLink.className = 'dynamic-admin-link text-amber-400 hover:text-amber-300 font-bold transition-colors text-sm px-2 py-1';
                adminDesktopLink.innerHTML = '⚙️ Inserisci';
                // Lo inseriamo subito prima delle info utente o del logout
                userInfo.parentNode.insertBefore(adminDesktopLink, userInfo);
            }

            // 2. Inserimento per Mobile (dentro il blocco mobile-auth-user se presente)
            if (mobileUser && mobileUser.parentNode) {
                const adminMobileLink = document.createElement('a');
                adminMobileLink.href = '/admin';
                adminMobileLink.className = 'dynamic-admin-link block text-amber-400 hover:text-amber-300 font-bold py-2';
                adminMobileLink.innerHTML = '⚙️ Inserisci';
                mobileUser.appendChild(adminMobileLink);
            }
        }

    } else {
        loginBtn?.classList.remove('hidden');
        registerBtn?.classList.remove('hidden');
        userInfo?.classList.add('hidden');
        logoutBtn?.classList.add('hidden');

        mobileGuest?.classList.remove('hidden');
        mobileUser?.classList.add('hidden');
    }
}

async function updateNavbarAuthState() {
    const user = await fetchCurrentUser();
    renderNavbarAuthState(user);
}

function setupLogoutHandlers() {
    const logoutButtons = [
        document.getElementById('nav-logout-btn'),
        document.getElementById('mobile-logout-btn')
    ];

    logoutButtons.forEach((button) => {
        if (!button) return;

        button.addEventListener('click', async (e) => {
            e.preventDefault();

            const raw = localStorage.getItem('user');
            let userId = null;

            if (raw) {
                try {
                    userId = JSON.parse(raw).user_id;
                } catch {
                    userId = null;
                }
            }

            try {
                if (userId) {
                    await fetch('/api/logout', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        credentials: 'include',
                        body: JSON.stringify({ user_id: userId })
                    });
                }
            } catch {
                // Proseguiamo comunque con la pulizia locale
            }

            localStorage.removeItem('user');
            window.location.href = '/';
        });
    });
}
