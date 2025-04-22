/**
 * Force styles and handle missing resources
 * This script ensures consistent styling even when some resources are missing
 * Respects tabler.min.css as base styling
 */
document.addEventListener('DOMContentLoaded', function() {
    // Забезпечуємо завантаження tabler.min.css першим
    function ensureTablerFirst() {
        const links = document.querySelectorAll('link[rel="stylesheet"]');
        const tablerLink = Array.from(links).find(link => link.href.includes('tabler.min.css'));
        const head = document.head;
        
        if (tablerLink && head.firstChild !== tablerLink) {
            head.insertBefore(tablerLink, head.firstChild);
        }
    }
    
    // Додаємо параметр версії до CSS файлів для обходу кешу
    function addRandomQueryToCSS() {
        const links = document.querySelectorAll('link[rel="stylesheet"]');
        const now = new Date().getTime();
        
        links.forEach(link => {
            if (link.href.includes('main.css')) {
                // Додаємо випадковий параметр, щоб уникнути кешування
                const url = new URL(link.href);
                url.searchParams.set('v', now);
                link.href = url.toString();
            }
        });
    }
    
    // Примусово застосовуємо кастомні стилі для кнопок на основі tabler 
    function applyCustomButtonStyles() {
        // Primary кнопки
        document.querySelectorAll('.btn-primary').forEach(btn => {
            btn.style.backgroundColor = '#D5BDAF';
            btn.style.borderColor = '#D5BDAF';
        });
        
        // Outline Primary кнопки
        document.querySelectorAll('.btn-outline-primary').forEach(btn => {
            btn.style.color = '#D5BDAF';
            btn.style.borderColor = '#D5BDAF';
        });
    }
    
    // Переконуємося, що класи шрифтів застосовані
    function ensureFontClasses() {
        document.documentElement.classList.add('font-inter');
    }
    
    // Примусово застосовуємо шрифт Inter через CSS змінні tabler
    function applyInterFontThroughTablerVars() {
        document.documentElement.style.setProperty('--tblr-font-sans-serif', 
            '"Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif');
    }

    // Викликаємо функції
    ensureTablerFirst();
    addRandomQueryToCSS();
    setTimeout(applyCustomButtonStyles, 100); // Невелика затримка для завантаження стилів
    ensureFontClasses();
    applyInterFontThroughTablerVars();

    // Force the footer background color
    document.querySelectorAll('.color-footer-bg').forEach(element => {
        element.style.backgroundColor = '#D5BDAF';
    });

    // Handle links that point to non-existent files
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                window.scrollTo({
                    top: target.offsetTop - 80,
                    behavior: 'smooth'
                });
            }
        });
    });
}); 