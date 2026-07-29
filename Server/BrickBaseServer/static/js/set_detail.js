document.addEventListener("DOMContentLoaded", async () => {
    
    const urlParams = new URLSearchParams(window.location.search);
    const setNum = urlParams.get('set_num');

    const titleElement = document.getElementById('set-title');
    const container = document.getElementById("set-parts-container");

    if (!setNum) {
        titleElement.innerText = "Errore: Nessun Set selezionato";
        container.innerHTML = "";
        return;
    }

    titleElement.innerHTML = `Inventario del Set: <span class="text-amber-500">${setNum}</span>`;

    try {
        
        const response = await fetch(`/api/sets/${setNum}/parts`);
        
        if (!response.ok) throw new Error("Errore nel recupero dati");
        
        const parts = await response.json();
        container.innerHTML = ""; 

        if (!parts || parts.length === 0) {
            container.innerHTML = `<p class="text-slate-400 col-span-full">Questo set non ha pezzi registrati o l'inventario è vuoto.</p>`;
            return;
        }

        
        parts.forEach(part => {
            
            const imgUrl = part.part_img_url && part.part_img_url !== "NOT_FOUND"
                ? part.part_img_url
                : '/static/images/notfound.png'; 
            
            
            const colorBox = part.color_rgb
                ? `<div class="w-3 h-3 rounded-full border border-slate-500" style="background-color: #${part.color_rgb};"></div>`
                : '';

            

            const card = `
                <div class="relative bg-slate-900 border border-slate-800 rounded-xl p-4 flex flex-col items-center hover:border-slate-600 transition shadow-lg">
                    <!-- CORREZIONE: Rimosso il badge da sopra l'immagine -->
                    <img src="${imgUrl}" alt="${part.part_name || part.part_num}" class="h-20 object-contain mb-4 drop-shadow-md">
                    
                    <h3 class="font-bold text-slate-200 text-center text-xs mb-1 line-clamp-2" title="${part.part_name}">
                        ${part.part_name || 'Nome sconosciuto'}
                    </h3>
                    
                    <p class="text-xs text-slate-500 mb-3">${part.part_num}</p>
                    
                    <div class="flex items-center gap-1.5 mb-3">
                        ${colorBox}
                        <span class="text-xs text-slate-400 truncate max-w-[100px]">${part.color_name || 'N/D'}</span>
                    </div>
                    
                    <div class="mt-auto w-full flex justify-between items-center border-t border-slate-800 pt-3">
                        <span class="text-xs text-slate-400 uppercase tracking-wider font-semibold">Qtà</span>
                        <span class="font-bold text-amber-500 text-lg">${part.quantity}</span>
                </div>
            `;
            
            container.innerHTML += card;
        });

    } catch (error) {
        console.error("Errore:", error);
        container.innerHTML = `<p class="text-red-500 col-span-full">Si è verificato un errore durante il caricamento dell'inventario.</p>`;
    }
});