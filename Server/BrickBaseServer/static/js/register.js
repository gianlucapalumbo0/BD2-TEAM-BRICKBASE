document.getElementById('register-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const errorMsg = document.getElementById('error-msg');
            const successMsg = document.getElementById('success-msg');
            errorMsg.classList.add('hidden');
            
            const data = {
                first_name: document.getElementById('first_name').value,
                last_name: document.getElementById('last_name').value,
                email: document.getElementById('email').value,
                password: document.getElementById('password').value,
                role: document.getElementById('role').value
            };

            try {
                const response = await fetch('/api/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data)
                });

                if (response.ok) {
                    successMsg.classList.remove('hidden');
                    setTimeout(() => window.location.href = '/login', 2000); // Vai al login
                } else {
                    const resData = await response.json();
                    errorMsg.textContent = resData.error || "Errore durante la registrazione";
                    errorMsg.classList.remove('hidden');
                }
            } catch (err) {
                errorMsg.textContent = "Errore di connessione al server.";
                errorMsg.classList.remove('hidden');
            }
        });