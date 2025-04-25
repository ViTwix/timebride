/**
 * Inter Font Loader
 * Simple script to ensure Inter font is properly loaded and applied
 */
document.addEventListener('DOMContentLoaded', function() {
    // Add font-inter class to html and body
    document.documentElement.classList.add('font-inter');
    document.body.classList.add('font-inter');
    
    // Check if fonts are loaded
    if ('fonts' in document) {
        document.fonts.ready.then(function() {
            // Add fonts-loaded class once fonts are loaded
            document.documentElement.classList.add('fonts-loaded');
        });
    } else {
        // Fallback for browsers that don't support fonts API
        document.documentElement.classList.add('fonts-loaded');
    }
}); 