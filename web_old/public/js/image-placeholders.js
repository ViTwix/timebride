/**
 * Image Placeholders
 * This script preloads and handles images efficiently
 */
document.addEventListener('DOMContentLoaded', function() {
    // Preload all images as soon as the DOM loads
    const imagesToPreload = [
        '/img/calendar-feature.png',
        '/img/client-feature.png',
        '/img/contract-feature.png',
        '/img/photography-bg.jpg',
        '/img/avatars/avatar-1.jpg',
        '/img/avatars/avatar-2.jpg',
        '/img/avatars/avatar-3.jpg'
    ];

    // Preload all critical images immediately
    imagesToPreload.forEach(imgPath => {
        const img = new Image();
        img.src = imgPath;
    });

    // Create base64 placeholder for failed images
    function createPlaceholder(width, height, text) {
        const canvas = document.createElement('canvas');
        canvas.width = width;
        canvas.height = height;
        const ctx = canvas.getContext('2d');
        
        // Background color
        ctx.fillStyle = '#D5BDAF'; // Updated to secondary color
        ctx.fillRect(0, 0, width, height);
        
        // Text
        ctx.fillStyle = '#6D4C41';
        ctx.font = 'bold 14px Inter, sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(text, width/2, height/2);
        
        return canvas.toDataURL();
    }
    
    // Apply fallback for avatars
    const avatars = document.querySelectorAll('img[src*="avatar"]');
    avatars.forEach(avatar => {
        // Create backup of the original source
        const originalSrc = avatar.getAttribute('src');
        
        // Pre-create fallback in case image fails
        const fallbackSrc = createPlaceholder(96, 96, 'Аватар');
        
        // Apply a specific handling for avatars
        avatar.onerror = function() {
            if (!this.getAttribute('data-fallback-applied')) {
                this.setAttribute('data-fallback-applied', 'true');
                this.src = fallbackSrc;
            }
        };
        
        // Force load to catch any existing errors
        const tempImg = new Image();
        tempImg.onload = function() {
            // If original loads successfully, ensure it's set
            avatar.src = originalSrc;
        };
        tempImg.onerror = function() {
            // If original fails, apply fallback immediately
            avatar.src = fallbackSrc;
            avatar.setAttribute('data-fallback-applied', 'true');
        };
        tempImg.src = originalSrc;
    });

    // Handle all other images
    document.querySelectorAll('img:not([src*="avatar"])').forEach(img => {
        if (!img.getAttribute('data-placeholder-applied')) {
            img.addEventListener('error', function() {
                if (!this.getAttribute('data-placeholder-applied')) {
                    this.setAttribute('data-placeholder-applied', 'true');
                    let width = this.width || 300;
                    let height = this.height || 150;
                    let text = 'Зображення';
                    
                    // Extract name from path for better placeholder text
                    const pathParts = this.src.split('/');
                    const filename = pathParts[pathParts.length - 1].split('.')[0];
                    if (filename) {
                        text = filename.replace(/-/g, ' ');
                        // Capitalize first letter
                        text = text.charAt(0).toUpperCase() + text.slice(1);
                    }
                    
                    this.src = createPlaceholder(width, height, text);
                }
            });
        }
    });
}); 