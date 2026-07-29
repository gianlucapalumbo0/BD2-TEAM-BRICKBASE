// Variabile globale per salvare in memoria i set già recensiti
let reviewedSets = new Set();

document.addEventListener('DOMContentLoaded', async () => {
    const container = document.getElementById('all-sets-container');
    if (container) {
        // 1. Aspettiamo che vengano scaricati i set già recensiti prima di procedere
        await fetchReviewedSets();
        // 2. Adesso carichiamo tutti i set in modo sicuro con i dati aggiornati in memoria
        fetchAllSets();
    }
});

let currentPage = 1;
let isLoading = false;
let hasMore = true; 

// Funzione di supporto per recuperare i dati utente dal localStorage
function getUserData() {
    try {
        const rawUser = localStorage.getItem('user');
        return rawUser ? JSON.parse(rawUser) : null;
    } catch (e) {
        console.error("Errore nel parsing dei dati utente:", e);
        return null;
    }
}

async function fetchReviewedSets() {
    const userData = getUserData();
    if (!userData) {
        console.log("⚠️ Utente non loggato.");
        return;
    }

    try {
        console.log("🚀 Richiedo i set recensiti al server tramite cookie...");
        // CORRETTO: rimosso '=awatt ='
        const response = await fetch('/api/users/me/reviews', {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            const data = await response.json();
            console.log("📥 Set recensiti ricevuti dal backend:", data);
            
            if (Array.isArray(data)) {
                reviewedSets = new Set(data.map(num => String(num).trim())); 
                console.log("💾 Memoria 'reviewedSets' aggiornata:", Array.from(reviewedSets));
            }
        } else {
            console.error("❌ Errore dal server. Status:", response.status);
        }
    } catch (error) {
        console.error("❌ Errore di connessione:", error);
    }
}

async function fetchAllSets() {
    if (isLoading || !hasMore) return;
    
    isLoading = true;
    try {
        const response = await fetch(`/api/sets?page=${currentPage}&limit=20`);
        if (!response.ok) throw new Error('Errore nel recupero dei set');
        
        const sets = await response.json();
        const container = document.getElementById('all-sets-container');
        
        if (!sets || sets.length === 0) {
            hasMore = false; 
            if (currentPage === 1) {
                container.innerHTML = '<p class="text-slate-400 col-span-full">Nessun set trovato.</p>';
            }
            return;
        }

        if (currentPage === 1) {
            container.innerHTML = '';
        }

        const userData = getUserData();
        const isLoggedIn = !!userData;
        const userRole = userData?.role ? userData.role.toUpperCase() : null;

        sets.forEach(set => {
            // CORRETTO: cleanSetNum ora è definita all'inizio del ciclo per tutti gli utenti
            const cleanSetNum = String(set.set_num).trim();
            const detailUrl = `/set-detail?set_num=${cleanSetNum}`;

            // Gestione pulsante recensione
            let reviewButtonHTML = '';
            
            if (isLoggedIn && userRole === 'USER') {
                if (reviewedSets.has(cleanSetNum)) {
                    reviewButtonHTML = `
                        <button disabled class="mt-2 block w-full py-2 text-center bg-slate-700 text-slate-400 rounded-lg text-sm font-bold cursor-not-allowed">
                            Già recensito
                        </button>
                    `;
                } else {
                    const safeName = set.name.replace(/'/g, "\\'").replace(/"/g, '&quot;');
                    reviewButtonHTML = `
                        <button id="btn-review-${cleanSetNum}" onclick="openReviewModal('${cleanSetNum}', '${safeName}')" class="mt-2 block w-full py-2 text-center bg-blue-600 hover:bg-blue-500 text-white rounded-lg text-sm font-bold transition-all">
                            Aggiungi recensione
                        </button>
                    `;
                }
            }

            // --- RATING FORMATTATO PER RIGA UNICA ---
            const ratingValue = set.review_rating || set.ReviewRating || 0;
            let ratingHTML = '';
            
            if (ratingValue > 0) {
                ratingHTML = `
                    <span class="flex items-center gap-1 text-amber-400 font-bold">
                        ⭐ ${Number(ratingValue).toFixed(1)} / 5
                    </span>
                `;
            } else {
                ratingHTML = `<span class="text-slate-500 italic">Nessun voto</span>`;
            }

            // --- CARD HTML CON RIGA UNICA ---
            const card = `
                <div class="bg-slate-900 p-4 rounded-2xl border border-slate-800 overflow-hidden hover:border-amber-500/50 transition-all flex flex-col">
                    <div class="h-40 bg-slate-800 rounded-xl mb-4 flex items-center justify-center overflow-hidden">
                        <img src="https://cdn.rebrickable.com/media/sets/${set.set_num}.jpg" 
                            alt="${set.name}" 
                            class="w-full h-full object-contain p-2"
                            onerror="this.onerror=null; this.src='/static/images/notfound.png';">
                    </div>
                    
                    <h3 class="font-bold text-white truncate mb-1">${set.name}</h3>
                    
                    <div class="flex items-center justify-between text-sm mb-3">
                        <span class="text-slate-400">Codice: ${set.set_num}</span>
                        <div id="rating-box-${cleanSetNum}">
                            ${ratingHTML}
                        </div>
                    </div>
                    
                    <div class="mt-auto">
                        <a href="${detailUrl}" class="block w-full py-2 text-center bg-amber-500 hover:bg-amber-600 text-slate-950 rounded-lg text-sm font-bold transition-all">Vedi dettagli</a>
                        ${reviewButtonHTML}
                    </div>
                </div>
            `;
            container.innerHTML += card;
        });

        currentPage++; 
    } catch (error) {
        console.error('Errore:', error);
        if (currentPage === 1) {
            document.getElementById('all-sets-container').innerHTML = 
                '<p class="text-red-400 col-span-full">Errore durante il caricamento dei set.</p>';
        }
    } finally {
        isLoading = false;
    }
}

// ASCOLTATORE DI SCROLL
window.addEventListener('scroll', () => {
    const container = document.getElementById('all-sets-container');
    if (!container) return;

    const windowHeight = window.innerHeight;
    const scrollY = Math.ceil(window.scrollY); 
    const documentHeight = document.documentElement.scrollHeight;

    if ((windowHeight + scrollY) >= (documentHeight - 500)) {
        fetchAllSets();
    }
});

/* ----------------------------------------------------
   LOGICA MODALE E INVIO RECENSIONE
------------------------------------------------------ */

function openReviewModal(setNum, setName) {
    const modal = document.getElementById('review-modal');
    if (!modal) return;
    
    modal.style.display = 'flex';
    
    document.getElementById('review-set-num').value = setNum;
    document.getElementById('modal-set-name').innerText = `Set: ${setName} (${setNum})`;
    
    document.getElementById('review-text').value = '';
    const msgEl = document.getElementById('review-message');
    msgEl.classList.add('hidden');
    msgEl.innerText = '';
    
    document.getElementById('submit-review-btn').disabled = false;
}

function closeReviewModal() {
    const modal = document.getElementById('review-modal');
    if (modal) {
        modal.style.display = 'none';
    }
}

async function submitReview() {
    const setNum = document.getElementById('review-set-num').value;
    const reviewText = document.getElementById('review-text').value.trim();
    const msgEl = document.getElementById('review-message');
    const submitBtn = document.getElementById('submit-review-btn');

    if (!reviewText) {
        msgEl.innerText = "Non puoi inviare una recensione vuota.";
        msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
        return;
    }

    submitBtn.disabled = true;
    msgEl.innerText = "Analisi della recensione con l'AI in corso...";
    msgEl.className = "text-sm font-semibold mb-4 text-blue-400 block";

    try {
        const response = await fetch(`/api/sets/${setNum}/review`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({ review: reviewText })
        });

        const data = await response.json();

        if (response.ok) {
            msgEl.innerText = "Recensione aggiunta con successo!";
            msgEl.className = "text-sm font-semibold mb-4 text-green-500 block";
            
            reviewedSets.add(setNum);

            // 1. Aggiorna il pulsante in "Già recensito"
            const cardButton = document.getElementById(`btn-review-${setNum}`);
            if (cardButton) {
                cardButton.outerHTML = `
                    <button disabled class="mt-2 block w-full py-2 text-center bg-slate-700 text-slate-400 rounded-lg text-sm font-bold cursor-not-allowed">
                        Già recensito
                    </button>
                `;
            }

            // 2. Aggiorna la media voti sulla card in tempo reale
            const ratingBox = document.getElementById(`rating-box-${setNum}`);
            if (ratingBox && data.review_rating) {
                ratingBox.innerHTML = `
                    <span class="flex items-center gap-1 text-amber-400 font-bold">
                        ⭐ ${Number(data.review_rating).toFixed(1)} / 5
                    </span>
                `;
            }

            setTimeout(() => {
                closeReviewModal();
            }, 2000);

        } else if (response.status === 409) {
            msgEl.innerText = "Hai già recensito questo set! Non è possibile modificarla.";
            msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
            
            reviewedSets.add(setNum);
            const cardButton = document.getElementById(`btn-review-${setNum}`);
            if (cardButton) {
                cardButton.outerHTML = `<button disabled class="mt-2 block w-full py-2 text-center bg-slate-700 text-slate-400 rounded-lg text-sm font-bold cursor-not-allowed">Già recensito</button>`;
            }

        } else {
            msgEl.innerText = data.error || "Errore durante l'inserimento della recensione.";
            msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
            submitBtn.disabled = false;
        }
    } catch (error) {
        console.error("Errore invio recensione:", error);
        msgEl.innerText = "Errore di connessione. Riprova più tardi.";
        msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
        submitBtn.disabled = false;
    }
}