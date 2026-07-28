document.getElementById('login-form').addEventListener('submit', async (e) => {
            e.preventDefault();

            const errorMsg = document.getElementById('error-msg');
            errorMsg.classList.add('hidden');

            const data = {
                email: document.getElementById('email').value,
                password: document.getElementById('password').value
            };

            try {
                const response = await fetch('/api/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    credentials: 'include',
                    body: JSON.stringify(data)
                });

                if (response.ok) {
                    // Il cookie è stato impostato dal server. Reindirizziamo l'utente alla home.
                    const userData = await response.json();

                    // Salviamo il nome e l'id dell'utente nel localStorage per poterlo mostrare nella navbar!
                    localStorage.setItem('user', JSON.stringify(userData));

                    window.location.href = '/';
                } else {
                    const resData = await response.json();
                    errorMsg.textContent = resData.error || "Email o password errati";
                    errorMsg.classList.remove('hidden');
                }
            } catch (err) {
                errorMsg.textContent = "Errore di connessione al server.";
                errorMsg.classList.remove('hidden');
            }
        });