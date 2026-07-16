document.addEventListener('DOMContentLoaded', () => {
    const container = document.getElementById('all-parts-container');
    if (container) {
        fetchAllParts();
    }
});

let currentPage = 1;
let isLoading = false;
let hasMore = true; 

async function fetchAllParts() {
    if (isLoading || !hasMore) return;
    
    isLoading = true;
    try {
        
        const response = await fetch(`/api/parts?page=${currentPage}&limit=20`);
        if (!response.ok) throw new Error('Errore nel recupero dei pezzi');
        
        const parts = await response.json();
        const container = document.getElementById('all-parts-container');
        
        
        if (!parts || parts.length === 0) {
            hasMore = false; 
            
            
            if (currentPage === 1) {
                container.innerHTML = '<p class="text-slate-400 col-span-full">Nessun pezzo trovato.</p>';
            }
            return;
        }

        
        if (currentPage === 1) {
            container.innerHTML = '';
        }

        parts.forEach(part => {
            // Se nel database c'è l'URL usiamo quello, altrimenti il notfound.
            const imageUrl = part.part_img_url ? part.part_img_url : '/static/images/notfound.png';

            const card = `
                <div class="bg-slate-900 p-4 rounded-2xl border border-slate-800 overflow-hidden hover:border-amber-500/50 transition-all flex flex-col justify-between h-full">
                    <div>
                        <div class="h-40 bg-slate-800 rounded-xl mb-4 flex items-center justify-center overflow-hidden">
                            <img src="${imageUrl}" 
                                alt="${part.name}" 
                                class="w-full h-full object-contain p-2"
                                onerror="this.onerror=null; this.src='/static/images/notfound.png';">
                        </div>
                        <h3 class="font-bold text-white truncate" title="${part.name}">${part.name}</h3>
                        <p class="text-slate-400 text-sm mb-4">Codice: ${part.part_num}</p>
                    </div>
                    <a href="/part/${part.part_num}" class="block w-full py-2 mt-4 text-center bg-slate-800 hover:bg-slate-700 rounded-lg text-sm font-semibold transition-all">Vedi dettagli</a>
                </div>
            `;
            container.innerHTML += card;
        });

        currentPage++; 
    } catch (error) {
        console.error('Errore:', error);
        
        if (currentPage === 1) {
            document.getElementById('all-parts-container').innerHTML = 
                '<p class="text-red-400 col-span-full">Errore durante il caricamento dei pezzi.</p>';
        }
    } finally {
        isLoading = false;
    }
}


window.addEventListener('scroll', () => {
    const container = document.getElementById('all-parts-container');
    if (!container) return;

    
    const windowHeight = window.innerHeight;
    const scrollY = Math.ceil(window.scrollY); 
    const documentHeight = document.documentElement.scrollHeight;

   
    if ((windowHeight + scrollY) >= (documentHeight - 500)) {
        fetchAllParts();
    }
});