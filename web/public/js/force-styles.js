/**
 * Force styles and handle missing resources
 * This script ensures consistent styling even when some resources are missing
 */
document.addEventListener('DOMContentLoaded', function() {
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
    
    // Примусово застосовуємо стилі для кнопок, якщо щось не працює
    function applyButtonStyles() {
        // Primary кнопки
        document.querySelectorAll('.btn-primary').forEach(btn => {
            btn.style.backgroundColor = '#D5BDAF';
            btn.style.borderColor = '#D5BDAF';
            btn.style.color = '#FFFFFF';
        });
        
        // Outline Primary кнопки
        document.querySelectorAll('.btn-outline-primary').forEach(btn => {
            btn.style.color = '#D5BDAF';
            btn.style.borderColor = '#D5BDAF';
            btn.style.backgroundColor = 'transparent';
        });
        
        // Outline Dark кнопки при наведенні - додаємо слухачі подій
        document.querySelectorAll('.btn-outline-dark').forEach(btn => {
            btn.addEventListener('mouseover', function() {
                this.style.backgroundColor = '#2A2A2A';
                this.style.borderColor = '#2A2A2A';
                this.style.color = '#FFFFFF';
            });
            
            btn.addEventListener('mouseout', function() {
                this.style.backgroundColor = 'transparent';
                this.style.borderColor = '#2A2A2A';
                this.style.color = '#2A2A2A';
            });
        });
    }
    
    // Переконуємося, що класи шрифтів застосовані
    function ensureFontClasses() {
        document.documentElement.classList.add('font-inter');
        document.body.classList.add('font-inter');
    }
    
    // Викликаємо функції
    addRandomQueryToCSS();
    setTimeout(applyButtonStyles, 100); // Невелика затримка для завантаження стилів
    ensureFontClasses();

    // Apply Inter font to all elements
    document.querySelectorAll('*').forEach(element => {
        if (window.getComputedStyle(element).fontFamily !== 'Inter') {
            element.style.fontFamily = 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif';
        }
    });

    // Примусово застосовуємо кольори для тексту
    document.querySelectorAll('.text-secondary').forEach(el => {
        el.style.color = '#4A4744 !important';
    });
    
    // Примусово застосовуємо кольори для заголовків
    document.querySelectorAll('h1.display-3, h2.display-5').forEach(el => {
        el.style.color = '#4A4744';
    });

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