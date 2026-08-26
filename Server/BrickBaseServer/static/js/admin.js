// 1. PROTEZIONE FRONTEND DELLA PAGINA
document.addEventListener('DOMContentLoaded', () => {
    const rawUser = localStorage.getItem('user');
    
    if (!rawUser) {
        window.location.href = '/';
        return;
    }

    try {
        const user = JSON.parse(rawUser);
        if (!user.role || user.role.toUpperCase() !== 'ADMIN') {
            alert("Accesso negato. Questa pagina è riservata agli amministratori.");
            window.location.href = '/';
        }
    } catch (e) {
        window.location.href = '/';
    }

    // Inizializza i listener per la ricerca nella sezione MODIFICA
    initUpdateSearchListeners();
});

// 2. LOGICA CAMBIO SEZIONI ADMIN (Tabs: Inserisci, Modifica, Cancella)
function showAdminTab(sectionName, clickedBtn) {
    document.querySelectorAll('.admin-section').forEach(sec => {
        sec.classList.add('hidden');
        sec.classList.remove('block');
    });
    
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('text-amber-500', 'border-b-2', 'border-amber-500');
        btn.classList.add('text-slate-400');
    });
    
    const target = document.getElementById(`section-${sectionName}`);
    if (target) {
        target.classList.remove('hidden');
        target.classList.add('block');
    }
    
    if (clickedBtn) {
        clickedBtn.classList.remove('text-slate-400');
        clickedBtn.classList.add('text-amber-500', 'border-b-2', 'border-amber-500');
    }
}

// 3. SET BUILDER: Array temporaneo e gestione inventario parti del set (INSERISCI)
let currentInventory = [];

function addPartToInventory() {
    const partNumInput = document.getElementById('input-part-num');
    const partNum = partNumInput.value.trim();
    const colorId = parseInt(document.getElementById('input-color-id').value);
    const quantity = parseInt(document.getElementById('input-quantity').value);
    const isSpare = document.getElementById('input-is-spare').checked;

    if (!partNum || isNaN(colorId) || isNaN(quantity) || quantity <= 0) {
        alert("Seleziona un pezzo valido tramite la barra di ricerca, seleziona un colore valido e inserisci una Quantità corretta.");
        return;
    }

    currentInventory.push({
        part_num: partNum,
        color_id: colorId,
        quantity: quantity,
        is_spare: isSpare,
        part_name: partNumInput.dataset.partName || "Pezzo LEGO",
        part_img_url: partNumInput.dataset.partImg || ""
    });

    document.getElementById('input-part-search').value = '';
    partNumInput.value = '';
    document.getElementById('input-color-search').value = '';
    document.getElementById('input-color-id').value = '';
    document.getElementById('input-quantity').value = '1';
    document.getElementById('input-is-spare').checked = false;
    
    const feedbackEl = document.getElementById('part-feedback');
    if (feedbackEl) {
        feedbackEl.classList.add('hidden');
        feedbackEl.innerHTML = '';
    }

    renderInventoryTable();
}

function removePartFromIndex(index) {
    currentInventory.splice(index, 1);
    renderInventoryTable();
}

function renderInventoryTable() {
    const tbody = document.getElementById('inventory-table-body');
    const emptyMsg = document.getElementById('empty-inventory-msg');
    if (!tbody) return;
    
    tbody.innerHTML = '';

    if (currentInventory.length === 0) {
        if (emptyMsg) emptyMsg.classList.remove('hidden');
        return;
    }

    if (emptyMsg) emptyMsg.classList.add('hidden');

    currentInventory.forEach((item, index) => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td class="py-3 px-3 flex items-center gap-3">
                <img src="${item.part_img_url || 'https://via.placeholder.com/32'}" class="w-8 h-8 object-contain bg-slate-950 rounded p-0.5 border border-slate-800">
                <div>
                    <div class="font-semibold text-slate-200 text-xs">${item.part_name}</div>
                    <div class="font-mono text-amber-400 text-xs">${item.part_num}</div>
                </div>
            </td>
            <td class="py-3 px-3">${item.color_id}</td>
            <td class="py-3 px-3">${item.quantity}</td>
            <td class="py-3 px-3">${item.is_spare ? 'Sì' : 'No'}</td>
            <td class="py-3 px-3 text-right">
                <button type="button" onclick="removePartFromIndex(${index})" class="text-red-400 hover:text-red-300 text-xs font-bold px-2 py-1 bg-red-950/40 rounded border border-red-900/50">Elimina</button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

// 4. LOGICA FORM "INSERISCI SET" (POST /api/addset)
async function submitAddSet(event) {
    event.preventDefault();

    const setNum = document.getElementById('set-num').value.trim();
    const setName = document.getElementById('set-name').value.trim();
    const setYear = parseInt(document.getElementById('set-year').value);
    const setThemeID = parseInt(document.getElementById('set-theme-id').value);
    const msgEl = document.getElementById('msg-set');

    if (currentInventory.length === 0) {
        msgEl.innerText = "Aggiungi almeno un pezzo all'inventario del set.";
        msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
        setTimeout(() => {
            msgEl.classList.add('hidden');
            msgEl.innerText = '';
        }, 4000);
        return;
    }

    const totalParts = currentInventory.reduce((acc, item) => acc + item.quantity, 0);

    const cleanInventory = currentInventory.map(item => ({
        part_num: item.part_num,
        color_id: item.color_id,
        quantity: item.quantity,
        is_spare: item.is_spare
    }));

    const payload = {
        set_num: setNum,
        name: setName,
        year: setYear,
        theme_id: setThemeID,
        num_parts: totalParts,
        parts_inventory: cleanInventory,
        user_reviews: [],
        review_rating: 0.0
    };

    try {
        const response = await fetch('/api/addset', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify(payload)
        });

        const data = await response.json();

        if (response.ok) {
            msgEl.innerText = "Set creato e salvato con successo nel database!";
            msgEl.className = "text-sm font-semibold mb-4 text-green-500 block";
            
            document.getElementById('form-add-set').reset();
            currentInventory = [];
            renderInventoryTable();

            ['set-num-feedback', 'set-name-feedback', 'set-year-feedback', 'theme-feedback', 'part-feedback'].forEach(id => {
                const el = document.getElementById(id);
                if (el) {
                    el.innerText = '';
                    el.classList.add('hidden');
                }
            });

            ['set-theme-id', 'input-part-num', 'input-color-id'].forEach(id => {
                const el = document.getElementById(id);
                if (el) el.value = '';
            });

            setTimeout(() => {
                msgEl.classList.add('hidden');
                msgEl.innerText = '';
            }, 5000);
        } else {
            msgEl.innerText = data.error || "Errore durante la creazione del set.";
            msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
            setTimeout(() => {
                msgEl.classList.add('hidden');
                msgEl.innerText = '';
            }, 6000);
        }
    } catch (error) {
        console.error("Errore:", error);
        msgEl.innerText = "Errore di connessione col server.";
        msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
    }
}

// 5. LOGICA FORM "MODIFICA SET" (PUT /api/set/:set_num)
let updateInventory = [];

// 1. CARICAMENTO DEL SET DAL DATABASE
async function fetchSetDetailsForUpdate() {
    const setNumInput = document.getElementById('update-set-num');
    const setNum = setNumInput ? setNumInput.value.trim() : '';
    const msgEl = document.getElementById('msg-load-set');
    const updateForm = document.getElementById('form-update-set');

    if (!setNum) {
        if (msgEl) {
            msgEl.innerText = "Inserisci un codice set valido.";
            msgEl.className = "text-xs mt-2 text-red-400 block font-medium";
            msgEl.classList.remove('hidden');
        }
        return;
    }

    try {
        const response = await fetch(`/api/set/${encodeURIComponent(setNum)}`, {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            const set = await response.json();
            
            document.getElementById('update-set-name').value = set.name || '';
            document.getElementById('update-set-year').value = set.year || '';
            document.getElementById('update-set-theme-id').value = set.theme_id || 0;

            updateInventory = set.parts_inventory ? [...set.parts_inventory] : [];
            renderUpdateInventoryTable();

            if (msgEl) {
                msgEl.innerText = `✔️ Set "${set.name}" caricato con successo!`;
                msgEl.className = "text-xs mt-2 text-green-400 block font-medium";
                msgEl.classList.remove('hidden');
            }
            if (updateForm) updateForm.classList.remove('hidden');
        } else {
            const data = await response.json();
            if (msgEl) {
                msgEl.innerText = data.error || "Set non trovato nel database.";
                msgEl.className = "text-xs mt-2 text-red-400 block font-medium";
                msgEl.classList.remove('hidden');
            }
            if (updateForm) updateForm.classList.add('hidden');
        }
    } catch (error) {
        console.error("Errore durante il recupero del set:", error);
        if (msgEl) {
            msgEl.innerText = "Errore di connessione col server.";
            msgEl.className = "text-xs mt-2 text-red-400 block font-medium";
            msgEl.classList.remove('hidden');
        }
    }
}

// 2. RENDERING TABELLA INVENTARIO MODIFICA
function renderUpdateInventoryTable() {
    const tbody = document.getElementById('update-inventory-table-body');
    const emptyMsg = document.getElementById('update-empty-inventory-msg');
    if (!tbody) return;

    tbody.innerHTML = '';

    if (updateInventory.length === 0) {
        if (emptyMsg) emptyMsg.classList.remove('hidden');
        return;
    }

    if (emptyMsg) emptyMsg.classList.add('hidden');

    updateInventory.forEach((item, index) => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td class="py-3 px-3">
                <div class="font-semibold text-slate-200 text-xs">${item.part_name || 'Pezzo LEGO'}</div>
                <div class="font-mono text-amber-400 text-xs">${item.part_num}</div>
            </td>
            <td class="py-3 px-3">${item.color_id}</td>
            <td class="py-3 px-3">${item.quantity}</td>
            <td class="py-3 px-3">${item.is_spare ? 'Sì' : 'No'}</td>
            <td class="py-3 px-3 text-right">
                <button type="button" onclick="removePartFromUpdateIndex(${index})" class="text-red-400 hover:text-red-300 text-xs font-bold px-2 py-1 bg-red-950/40 rounded border border-red-900/50">Elimina</button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

// 3. AGGIUNGERE UN PEZZO ALL'INVENTARIO MODIFICA
function addPartToUpdateInventory() {
    const partNumInput = document.getElementById('update-part-num');
    const partNum = partNumInput ? partNumInput.value.trim() : '';
    const colorId = parseInt(document.getElementById('update-color-id').value);
    const quantity = parseInt(document.getElementById('update-quantity').value);
    const isSpare = document.getElementById('update-is-spare').checked;

    if (!partNum || isNaN(colorId) || isNaN(quantity) || quantity <= 0) {
        alert("Seleziona un pezzo e un colore validi dal menu ed inserisci una quantità corretta.");
        return;
    }

    updateInventory.push({
        part_num: partNum,
        color_id: colorId,
        quantity: quantity,
        is_spare: isSpare,
        part_name: partNumInput.dataset.partName || "Pezzo LEGO"
    });

    document.getElementById('update-part-search').value = '';
    partNumInput.value = '';
    document.getElementById('update-color-search').value = '';
    document.getElementById('update-color-id').value = '';
    document.getElementById('update-quantity').value = '1';
    document.getElementById('update-is-spare').checked = false;

    renderUpdateInventoryTable();
}

// 4. RIMUOVERE UN PEZZO DALL'INVENTARIO MODIFICA
function removePartFromUpdateIndex(index) {
    updateInventory.splice(index, 1);
    renderUpdateInventoryTable();
}

// 5. INVIO FORM MODIFICA AL BACKEND (PUT /api/set/:set_num)
// 5. INVIO FORM MODIFICA AL BACKEND (PUT /api/set/:set_num)
async function submitUpdateSet(event) {
    event.preventDefault();

    const setNum = document.getElementById('update-set-num').value.trim();
    const setName = document.getElementById('update-set-name').value.trim();
    const setYear = parseInt(document.getElementById('update-set-year').value);
    const setThemeID = parseInt(document.getElementById('update-set-theme-id').value) || 0;
    const msgEl = document.getElementById('msg-update-set');

    const totalParts = updateInventory.reduce((acc, item) => acc + item.quantity, 0);

    const cleanInventory = updateInventory.map(item => ({
        part_num: item.part_num,
        color_id: item.color_id,
        quantity: item.quantity,
        is_spare: item.is_spare
    }));

    const payload = {
        name: setName,
        year: setYear,
        theme_id: setThemeID,
        num_parts: totalParts,
        parts_inventory: cleanInventory
    };

    try {
        const response = await fetch(`/api/set/${encodeURIComponent(setNum)}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify(payload)
        });

        const data = await response.json();

        if (response.ok) {
            // Messaggio di conferma
            msgEl.innerText = "Set e inventario aggiornati con successo!";
            msgEl.className = "text-sm font-semibold mb-4 text-green-500 block";
            msgEl.classList.remove('hidden');

            // 1. Resetta e nasconde il form principale di modifica
            const updateForm = document.getElementById('form-update-set');
            if (updateForm) {
                updateForm.reset();
                updateForm.classList.add('hidden');
            }

            // 2. Svuota il campo di ricerca del codice set e il relativo messaggio di caricamento
            const setNumInput = document.getElementById('update-set-num');
            if (setNumInput) setNumInput.value = '';

            const msgLoad = document.getElementById('msg-load-set');
            if (msgLoad) {
                msgLoad.innerText = '';
                msgLoad.classList.add('hidden');
            }

            // 3. Svuota l'array dell'inventario e ri-renderizza la tabella (che tornerà vuota)
            updateInventory = [];
            renderUpdateInventoryTable();

            // 4. Pulizia degli input di ricerca pezzi e colori interni
            ['update-part-search', 'update-part-num', 'update-color-search', 'update-color-id'].forEach(id => {
                const el = document.getElementById(id);
                if (el) el.value = '';
            });

            // 5. Fa scomparire il messaggio di successo dopo 4 secondi
            setTimeout(() => {
                msgEl.classList.add('hidden');
                msgEl.innerText = '';
            }, 4000);

        } else {
            msgEl.innerText = data.error || "Errore durante l'aggiornamento del set.";
            msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
            msgEl.classList.remove('hidden');
        }
    } catch (error) {
        console.error("Errore:", error);
        msgEl.innerText = "Errore di connessione col server.";
        msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
        msgEl.classList.remove('hidden');
    }
}

// 6. LOGICA FORM "CANCELLA SET" (DELETE /api/set/:set_num)
// 6. LOGICA FORM "CANCELLA SET" (DELETE /api/set/:set_num)
async function submitDeleteSet(event) {
    event.preventDefault();

    const setNumInput = document.getElementById('delete-set-num');
    const setNum = setNumInput ? setNumInput.value.trim() : '';
    const msgEl = document.getElementById('msg-delete-set');

    if (!setNum) {
        if (msgEl) {
            msgEl.innerText = "Inserisci un codice set valido.";
            msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
            msgEl.classList.remove('hidden');
        }
        return;
    }

    if (!confirm(`Sei sicuro di voler eliminare permanentemente il set ${setNum}?`)) {
        return;
    }

    try {
        const response = await fetch(`/api/set/${encodeURIComponent(setNum)}`, {
            method: 'DELETE',
            credentials: 'include'
        });

        const data = await response.json();

        if (response.ok) {
            // 1. Messaggio di conferma successo
            msgEl.innerText = `✔️ Set ${setNum} eliminato con successo dal database!`;
            msgEl.className = "text-sm font-semibold mb-4 text-green-500 block";
            msgEl.classList.remove('hidden');

            // 2. Resetta il form e pulisce l'input
            const deleteForm = document.getElementById('form-delete-set');
            if (deleteForm) deleteForm.reset();
            if (setNumInput) setNumInput.value = '';

            // 3. Nasconde il messaggio dopo 4 secondi per lasciare la pagina pulita
            setTimeout(() => {
                msgEl.classList.add('hidden');
                msgEl.innerText = '';
            }, 4000);

        } else {
            msgEl.innerText = data.error || "Errore durante l'eliminazione del set.";
            msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
            msgEl.classList.remove('hidden');
        }
    } catch (error) {
        console.error("Errore:", error);
        msgEl.innerText = "Errore di connessione col server.";
        msgEl.className = "text-sm font-semibold mb-4 text-red-500 block";
        msgEl.classList.remove('hidden');
    }
}

// CONTROLLO DISPONIBILITÀ CAMPO
async function checkFieldAvailability(field, inputId, feedbackId) {
    const inputEl = document.getElementById(inputId);
    const feedbackEl = document.getElementById(feedbackId);
    if (!inputEl) return;
    const value = inputEl.value.trim();

    if (!value) {
        if (feedbackEl) feedbackEl.classList.add('hidden');
        return;
    }

    try {
        const response = await fetch(`/api/set/check?field=${field}&value=${encodeURIComponent(value)}`, {
            method: 'GET',
            credentials: 'include'
        });

        if (response.ok) {
            const data = await response.json();
            if (!feedbackEl) return;

            if (field === 'set_num') {
                if (data.exists) {
                    feedbackEl.innerText = "⚠️ Questo codice set esiste già nel database!";
                    feedbackEl.className = "text-xs mt-1.5 text-red-400 block font-medium";
                } else {
                    feedbackEl.innerText = "✔️ Codice set disponibile.";
                    feedbackEl.className = "text-xs mt-1.5 text-green-400 block font-medium";
                }
            } else if (field === 'name') {
                if (data.exists) {
                    feedbackEl.innerText = "ℹ️ Nota: Esiste già un set con questo nome nel database.";
                    feedbackEl.className = "text-xs mt-1.5 text-amber-400 block font-medium";
                } else {
                    feedbackEl.classList.add('hidden');
                }
            }
        }
    } catch (error) {
        console.error(`Errore durante la verifica del campo ${field}:`, error);
    }
}

// CONTROLLO RANGE ANNO
function checkYearRange(inputEl) {
    const feedbackEl = document.getElementById('set-year-feedback');
    if (!inputEl || !feedbackEl) return;
    const year = parseInt(inputEl.value);
    const currentYear = new Date().getFullYear();
    const minYear = 1932;

    if (!inputEl.value) {
        feedbackEl.classList.add('hidden');
        return;
    }

    if (year < minYear || year > currentYear) {
        feedbackEl.innerText = `⚠️ L'anno deve essere compreso tra ${minYear} e ${currentYear}.`;
        feedbackEl.className = "text-xs mt-1.5 text-red-400 block font-medium";
    } else {
        feedbackEl.innerText = "✔️ Anno valido.";
        feedbackEl.className = "text-xs mt-1.5 text-green-400 block font-medium";
    }
}

// INIZIALIZZATORE LISTENER DI RICERCA PER SEZIONE MODIFICA
function initUpdateSearchListeners() {
    const updatePartSearch = document.getElementById('update-part-search');
    if (updatePartSearch) {
        let timeout;
        updatePartSearch.addEventListener('input', (e) => {
            clearTimeout(timeout);
            const query = e.target.value.trim();
            const suggestions = document.getElementById('update-part-suggestions');
            const partNumInput = document.getElementById('update-part-num');
            if (query.length < 2) {
                if (suggestions) suggestions.classList.add('hidden');
                return;
            }
            timeout = setTimeout(async () => {
                const res = await fetch(`/api/parts/search?q=${encodeURIComponent(query)}`);
                if (res.ok && suggestions) {
                    const parts = await res.json();
                    suggestions.innerHTML = '';
                    parts.forEach(p => {
                        const d = document.createElement('div');
                        d.className = "p-2.5 hover:bg-slate-800 cursor-pointer text-sm text-slate-200 border-b border-slate-800/80";
                        d.innerText = `${p.name} (${p.part_num})`;
                        d.onclick = () => {
                            updatePartSearch.value = p.part_num;
                            partNumInput.value = p.part_num;
                            partNumInput.dataset.partName = p.name;
                            suggestions.classList.add('hidden');
                        };
                        suggestions.appendChild(d);
                    });
                    suggestions.classList.remove('hidden');
                }
            }, 300);
        });
    }

    const updateColorSearch = document.getElementById('update-color-search');
    if (updateColorSearch) {
        let timeout;
        updateColorSearch.addEventListener('input', (e) => {
            clearTimeout(timeout);
            const query = e.target.value.trim();
            const suggestions = document.getElementById('update-color-suggestions');
            const colorIdInput = document.getElementById('update-color-id');
            if (query.length < 1) {
                if (suggestions) suggestions.classList.add('hidden');
                return;
            }
            timeout = setTimeout(async () => {
                const res = await fetch(`/api/colors/search?q=${encodeURIComponent(query)}`);
                if (res.ok && suggestions) {
                    const colors = await res.json();
                    suggestions.innerHTML = '';
                    colors.forEach(c => {
                        const d = document.createElement('div');
                        d.className = "p-2.5 hover:bg-slate-800 cursor-pointer text-sm text-slate-200 border-b border-slate-800/80";
                        d.innerText = `${c.name} (ID: ${c.color_id || c.id})`;
                        d.onclick = () => {
                            updateColorSearch.value = c.name;
                            colorIdInput.value = c.color_id || c.id;
                            suggestions.classList.add('hidden');
                        };
                        suggestions.appendChild(d);
                    });
                    suggestions.classList.remove('hidden');
                }
            }, 300);
        });
    }
}

// RICERCA TEMI
let themeSearchTimeout;
const themeSearchInput = document.getElementById('set-theme-search');
if (themeSearchInput) {
    themeSearchInput.addEventListener('input', function(e) {
        clearTimeout(themeSearchTimeout);
        const query = e.target.value.trim();
        const suggestionsEl = document.getElementById('theme-suggestions');
        const themeIdInput = document.getElementById('set-theme-id');
        const feedbackEl = document.getElementById('theme-feedback');

        if (themeIdInput) themeIdInput.value = '';

        if (query.length < 2) {
            if (suggestionsEl) suggestionsEl.classList.add('hidden');
            if (feedbackEl) feedbackEl.classList.add('hidden');
            return;
        }

        themeSearchTimeout = setTimeout(async () => {
            try {
                const response = await fetch(`/api/themes/search?q=${encodeURIComponent(query)}`, {
                    method: 'GET',
                    credentials: 'include'
                });

                if (response.ok && suggestionsEl) {
                    const themes = await response.json();
                    suggestionsEl.innerHTML = '';

                    if (themes.length === 0) {
                        suggestionsEl.classList.add('hidden');
                        if (feedbackEl) {
                            feedbackEl.innerText = "⚠️ Nessun tema trovato con questo nome.";
                            feedbackEl.className = "text-xs mt-1.5 text-amber-400 block font-medium";
                            feedbackEl.classList.remove('hidden');
                        }
                        return;
                    }

                    if (feedbackEl) feedbackEl.classList.add('hidden');

                    themes.forEach(theme => {
                        const div = document.createElement('div');
                        div.className = "p-3 hover:bg-slate-800 cursor-pointer text-sm text-slate-200 border-b border-slate-800/80 last:border-none transition-colors";
                        div.innerHTML = `<span>${theme.name}</span> <span class="text-xs text-slate-500 float-right">ID: ${theme.theme_id}</span>`;
                        
                        div.onclick = () => {
                            themeSearchInput.value = theme.name;
                            if (themeIdInput) themeIdInput.value = theme.theme_id;
                            suggestionsEl.classList.add('hidden');
                            
                            if (feedbackEl) {
                                feedbackEl.innerText = `✔️ Tema associato correttamente (ID: ${theme.theme_id})`;
                                feedbackEl.className = "text-xs mt-1.5 text-green-400 block font-medium";
                                feedbackEl.classList.remove('hidden');
                            }
                        };
                        suggestionsEl.appendChild(div);
                    });

                    suggestionsEl.classList.remove('hidden');
                }
            } catch (error) {
                console.error("Errore durante la ricerca dei temi:", error);
            }
        }, 300);
    });
}

// RICERCA PEZZI INSERISCI
let partSearchTimeout;
const inputPartSearch = document.getElementById('input-part-search');
if (inputPartSearch) {
    inputPartSearch.addEventListener('input', function(e) {
        clearTimeout(partSearchTimeout);
        const query = e.target.value.trim();
        const suggestionsEl = document.getElementById('part-suggestions');
        const partNumInput = document.getElementById('input-part-num');
        const feedbackEl = document.getElementById('part-feedback');

        if (partNumInput) partNumInput.value = '';

        if (query.length < 2) {
            if (suggestionsEl) suggestionsEl.classList.add('hidden');
            if (feedbackEl) feedbackEl.classList.add('hidden');
            return;
        }

        partSearchTimeout = setTimeout(async () => {
            try {
                const response = await fetch(`/api/parts/search?q=${encodeURIComponent(query)}`, {
                    method: 'GET',
                    credentials: 'include'
                });

                if (response.ok && suggestionsEl) {
                    const parts = await response.json();
                    suggestionsEl.innerHTML = '';

                    if (parts.length === 0) {
                        suggestionsEl.classList.add('hidden');
                        if (feedbackEl) {
                            feedbackEl.innerText = "⚠️ Nessun pezzo trovato.";
                            feedbackEl.className = "text-xs mt-1.5 text-amber-400 block font-medium";
                            feedbackEl.classList.remove('hidden');
                        }
                        return;
                    }

                    if (feedbackEl) feedbackEl.classList.add('hidden');

                    parts.forEach(part => {
                        const div = document.createElement('div');
                        div.className = "flex items-center gap-3 p-2.5 hover:bg-slate-800 cursor-pointer text-sm text-slate-200 border-b border-slate-800/80 last:border-none transition-colors";
                        
                        const imgSrc = (part.part_img_url && part.part_img_url !== "NOT_FOUND") 
                            ? part.part_img_url 
                            : 'https://via.placeholder.com/40?text=No+Img';

                        div.innerHTML = `
                            <img src="${imgSrc}" alt="${part.name}" class="w-10 h-10 object-contain bg-slate-950 rounded p-1 border border-slate-700">
                            <div class="flex-1 overflow-hidden">
                                <div class="font-bold truncate text-slate-200">${part.name}</div>
                                <div class="text-xs text-amber-400 font-mono">Codice: ${part.part_num}</div>
                            </div>
                        `;
                        
                        div.onclick = () => {
                            inputPartSearch.value = part.part_num;
                            if (partNumInput) {
                                partNumInput.value = part.part_num;
                                partNumInput.dataset.partName = part.name;
                                partNumInput.dataset.partImg = imgSrc;
                            }

                            suggestionsEl.classList.add('hidden');
                            
                            if (feedbackEl) {
                                feedbackEl.innerHTML = `
                                    <div class="flex items-center gap-3 p-3 bg-slate-900 border border-slate-700 rounded-xl shadow-md">
                                        <img src="${imgSrc}" class="w-10 h-10 object-contain bg-slate-950 rounded p-1 border border-slate-800">
                                        <div>
                                            <div class="text-xs font-bold text-slate-200">${part.name}</div>
                                            <div class="text-xs text-amber-400 font-mono">Codice: ${part.part_num} ✔️ Selezionato con successo</div>
                                        </div>
                                    </div>
                                `;
                                feedbackEl.classList.remove('hidden');
                            }
                        };

                        suggestionsEl.appendChild(div);
                    });

                    suggestionsEl.classList.remove('hidden');
                }
            } catch (error) {
                console.error("Errore durante la ricerca dei pezzi:", error);
            }
        }, 300);
    });
}

// RICERCA COLORI INSERISCI
let colorSearchTimeout;
const inputColorSearch = document.getElementById('input-color-search');
if (inputColorSearch) {
    inputColorSearch.addEventListener('input', function(e) {
        clearTimeout(colorSearchTimeout);
        const query = e.target.value.trim();
        const suggestionsEl = document.getElementById('color-suggestions');
        const colorIdInput = document.getElementById('input-color-id');

        if (colorIdInput) colorIdInput.value = '';

        if (query.length < 1) {
            if (suggestionsEl) suggestionsEl.classList.add('hidden');
            return;
        }

        colorSearchTimeout = setTimeout(async () => {
            try {
                const response = await fetch(`/api/colors/search?q=${encodeURIComponent(query)}`, {
                    method: 'GET',
                    credentials: 'include'
                });

                if (response.ok && suggestionsEl) {
                    const colors = await response.json();
                    suggestionsEl.innerHTML = '';

                    if (colors.length === 0) {
                        suggestionsEl.classList.add('hidden');
                        return;
                    }

                    colors.forEach(color => {
                        const div = document.createElement('div');
                        div.className = "flex items-center gap-2.5 p-2.5 hover:bg-slate-800 cursor-pointer text-sm text-slate-200 border-b border-slate-800/80 last:border-none transition-colors";
                        
                        const rgb = color.rgb ? `#${color.rgb}` : 'transparent';
                        const colorIndicator = color.rgb ? `<span class="w-4 h-4 rounded-full border border-slate-600 inline-block shadow-inner" style="background-color: ${rgb}"></span>` : '';

                        div.innerHTML = `
                            ${colorIndicator}
                            <div class="flex-1 truncate">
                                <span class="font-medium text-slate-200">${color.name}</span>
                                <span class="text-xs text-slate-500 ml-1">(ID: ${color.color_id || color.id})</span>
                            </div>
                        `;
                        
                        div.onclick = () => {
                            inputColorSearch.value = color.name;
                            if (colorIdInput) colorIdInput.value = color.color_id || color.id;
                            suggestionsEl.classList.add('hidden');
                        };

                        suggestionsEl.appendChild(div);
                    });

                    suggestionsEl.classList.remove('hidden');
                }
            } catch (error) {
                console.error("Errore durante la ricerca dei colori:", error);
            }
        }, 300);
    });
}

// ESPORTAZIONE GLOBALE PER L'HTML
window.showAdminTab = showAdminTab;
window.addPartToInventory = addPartToInventory;
window.removePartFromIndex = removePartFromIndex;
window.submitAddSet = submitAddSet;
window.fetchSetDetailsForUpdate = fetchSetDetailsForUpdate;
window.addPartToUpdateInventory = addPartToUpdateInventory;
window.removePartFromUpdateIndex = removePartFromUpdateIndex;
window.submitUpdateSet = submitUpdateSet;
window.submitDeleteSet = submitDeleteSet;
window.checkFieldAvailability = checkFieldAvailability;
window.checkYearRange = checkYearRange;