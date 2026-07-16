document.addEventListener('DOMContentLoaded', () => {
   
    const container = document.getElementById('all-sets-container');
    if (container) {
        fetchAllSets();
    }
});

let currentPage = 1;
let isLoading = false;
let hasMore = true; // Per smettere di chiamare l'API se finiscono i set

async function fetchAllSets() {
    if (isLoading || !hasMore) return;
    
    isLoading = true;
    try {
        const response = await fetch(`/api/sets?page=${currentPage}&limit=20`);
        if (!response.ok) throw new Error('Errore nel recupero dei set');
        
        const sets = await response.json();
        const container = document.getElementById('all-sets-container');
        
        // GESTIONE NULL/VUOTO: Controlla che sets esista prima di leggerne la lunghezza
        if (!sets || sets.length === 0) {
            hasMore = false; 
            
            // Se è la prima pagina e non ci sono set, mostra un messaggio
            if (currentPage === 1) {
                container.innerHTML = '<p class="text-slate-400 col-span-full">Nessun set trovato.</p>';
            }
            return;
        }

        // RIMOZIONE TESTO: Se è la prima pagina, svuota il "Caricamento in corso..."
        if (currentPage === 1) {
            container.innerHTML = '';
        }

        sets.forEach(set => {
            const card = `
                <div class="bg-slate-900 p-4 rounded-2xl border border-slate-800 overflow-hidden hover:border-amber-500/50 transition-all">
                    <div class="h-40 bg-slate-800 rounded-xl mb-4 flex items-center justify-center overflow-hidden">
                        <img src="https://cdn.rebrickable.com/media/sets/${set.set_num}.jpg" 
                            alt="${set.name}" 
                            class="w-full h-full object-contain p-2"
                            onerror="this.onerror=null; this.src='/static/images/notfound.png';">
                    </div>
                    <h3 class="font-bold text-white truncate">${set.name}</h3>
                    <p class="text-slate-400 text-sm mb-4">Codice: ${set.set_num}</p>
                    <a href="/set/${set.set_num}" class="block w-full py-2 text-center bg-slate-800 hover:bg-slate-700 rounded-lg text-sm font-semibold transition-all">Vedi dettagli</a>
                </div>
            `;
            container.innerHTML += card;
        });

        currentPage++; 
    } catch (error) {
        console.error('Errore:', error);
        // In caso di errore alla prima pagina, diamo un feedback visivo
        if (currentPage === 1) {
            document.getElementById('all-sets-container').innerHTML = 
                '<p class="text-red-400 col-span-full">Errore durante il caricamento dei set.</p>';
        }
    } finally {
        isLoading = false;
    }
}

// ASCOLTATORE DI SCROLL MIGLIORATO
window.addEventListener('scroll', () => {
    const container = document.getElementById('all-sets-container');
    if (!container) return;

    // Raccogliamo le misure esatte cross-browser
    const windowHeight = window.innerHeight;
    const scrollY = Math.ceil(window.scrollY); // Arrotonda per difetto i pixel decimali
    const documentHeight = document.documentElement.scrollHeight;

    // Se ci troviamo a 500px (o meno) dalla fine della pagina, carica la prossima pagina
    if ((windowHeight + scrollY) >= (documentHeight - 500)) {
        fetchAllSets();
    }
});