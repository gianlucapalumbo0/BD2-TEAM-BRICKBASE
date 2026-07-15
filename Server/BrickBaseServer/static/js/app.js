// Aspetta che la pagina sia completamente caricata
document.addEventListener('DOMContentLoaded', () => {
    // Seleziona gli elementi della navbar
    const btn = document.getElementById('mobile-menu-btn');
    const menu = document.getElementById('mobile-menu');
    const iconOpen = document.getElementById('icon-open');
    const iconClose = document.getElementById('icon-close');

    // Al click sul pulsante mobile...
    if (btn) {
        btn.addEventListener('click', () => {
            // Mostra o nascondi il menu a tendina
            menu.classList.toggle('hidden');
            
            // Alterna le icone (hamburger <-> X)
            iconOpen.classList.toggle('hidden');
            iconClose.classList.toggle('hidden');
        });
    }
});

document.addEventListener('DOMContentLoaded', () => {
    fetchBestSets();
});

document.addEventListener('DOMContentLoaded', () => {
    fetchBestSets();
});

async function fetchBestSets() {
    try {
        // La rotta punta al controller GetBestSets definito in unprotected_routes.go[cite: 2]
        const response = await fetch('/api/bestsets'); 
        if (!response.ok) throw new Error('Errore nel recupero dei set');
        
        const sets = await response.json(); // Riceviamo l'array di modelli Set
        const container = document.getElementById('best-sets-container');
        
        if (!container) return; // Protezione se l'id non esiste nella pagina
        
        container.innerHTML = ''; // Pulisce il contenitore

        sets.forEach(set => {
            // Utilizziamo le chiavi definite nei tag json di models/set.go
            const card = `
                <div class="min-w-[280px] md:min-w-[300px] bg-slate-900 p-4 rounded-2xl border border-slate-800 snap-center">
                    <div class="h-40 bg-slate-800 rounded-xl mb-4 flex items-center justify-center text-4xl">🧱</div>
                    <h3 class="font-bold text-white truncate">${set.name}</h3>
                    <p class="text-slate-400 text-sm">Codice: ${set.set_num}</p>
                    <p class="text-amber-500 text-sm mb-4">Rating: ${set.review_rating ? set.review_rating.toFixed(1) : 'N/D'}</p>
                    <a href="/set/${set.set_num}" class="block w-full py-2 text-center bg-slate-800 hover:bg-slate-700 rounded-lg text-sm font-semibold transition-all">Dettagli</a>
                </div>
            `;
            container.innerHTML += card;
        });
    } catch (error) {
        console.error('Errore:', error);
    }
}