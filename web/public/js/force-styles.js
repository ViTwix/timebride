/**
 * TimeBride - Стилі та відображення
 */
document.addEventListener('DOMContentLoaded', function() {
    // Встановлення кольорової схеми
    function applyColorScheme() {
        // Основний колір
        const mainColor = getComputedStyle(document.documentElement).getPropertyValue('--color-sidebar-bg').trim();
        
        // Примусово застосовуємо до кнопок
        document.querySelectorAll('.btn-primary').forEach(btn => {
            btn.style.backgroundColor = mainColor;
            btn.style.borderColor = mainColor;
        });
        
        document.querySelectorAll('.btn-outline-primary').forEach(btn => {
            btn.style.color = mainColor;
            btn.style.borderColor = mainColor;
        });
    }
    
    // Переконуємося, що шрифти завантажені та застосовані
    function setupFonts() {
        document.documentElement.classList.add('font-inter');
    }
    
    // Виклик функцій
    applyColorScheme();
    setupFonts();
}); 