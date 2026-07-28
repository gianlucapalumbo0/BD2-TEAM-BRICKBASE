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
});
