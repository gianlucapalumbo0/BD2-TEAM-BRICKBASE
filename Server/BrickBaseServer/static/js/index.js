


document.addEventListener('DOMContentLoaded', () => {
    fetchBestSets();
});



async function fetchBestSets() {
    try {
        
        const response = await fetch('/api/bestsets'); 
        if (!response.ok) throw new Error('Errore nel recupero dei set');
        
        const sets = await response.json(); 
        const container = document.getElementById('best-sets-container');
        
        if (!container) return; 
        
        container.innerHTML = '';

        sets.forEach(set => {
            
            const card = `
                <div class="min-w-[280px] md:min-w-[300px] bg-slate-900 p-4 rounded-2xl border border-slate-800 snap-center overflow-hidden">
                    <div class="h-40 bg-slate-800 rounded-xl mb-4 flex items-center justify-center overflow-hidden">
                        <img src="https://cdn.rebrickable.com/media/sets/${set.set_num}.jpg" 
                            alt="${set.name}" 
                            class="w-full h-full object-contain p-2"
                            onerror="this.onerror=null; this.src='/static/images/notfound.png';">
                    </div>
                    
                    <h3 class="font-bold text-white truncate">${set.name}</h3>
                    <p class="text-slate-400 text-sm">Codice: ${set.set_num}</p>
                    <!-- ... resto della card ... -->
                </div>
            `;
            
            container.innerHTML += card;
        });
    } catch (error) {
        console.error('Errore:', error);
    }
}



