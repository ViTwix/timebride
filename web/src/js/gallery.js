// Gallery Component
class Gallery {
    constructor(container, options = {}) {
        this.container = container;
        this.options = {
            items: options.items || [],
            columns: options.columns || 3,
            gap: options.gap || 1,
            onItemClick: options.onItemClick || null,
            onLoadMore: options.onLoadMore || null,
            loading: false,
            hasMore: true
        };
        
        this.init();
    }

    init() {
        this.render();
        this.attachEventListeners();
    }

    render() {
        const gallery = document.createElement('div');
        gallery.className = 'gallery';
        gallery.style.setProperty('--columns', this.options.columns);
        gallery.style.setProperty('--gap', `${this.options.gap}rem`);
        
        // Gallery grid
        const grid = document.createElement('div');
        grid.className = 'gallery-grid';
        
        this.options.items.forEach(item => {
            const itemEl = this.createGalleryItem(item);
            grid.appendChild(itemEl);
        });
        
        gallery.appendChild(grid);
        
        // Load more button
        if (this.options.hasMore) {
            const loadMore = document.createElement('button');
            loadMore.className = 'gallery-load-more';
            loadMore.textContent = this.options.loading ? 'Loading...' : 'Load More';
            loadMore.disabled = this.options.loading;
            loadMore.onclick = () => this.options.onLoadMore?.();
            gallery.appendChild(loadMore);
        }
        
        this.container.innerHTML = '';
        this.container.appendChild(gallery);
    }

    createGalleryItem(item) {
        const itemEl = document.createElement('div');
        itemEl.className = 'gallery-item';
        
        // Image container
        const imgContainer = document.createElement('div');
        imgContainer.className = 'gallery-item-image';
        
        // Image
        const img = document.createElement('img');
        img.src = item.thumbnail_url || item.url;
        img.alt = item.title || '';
        img.loading = 'lazy';
        
        // Overlay
        const overlay = document.createElement('div');
        overlay.className = 'gallery-item-overlay';
        
        // Title
        const title = document.createElement('h3');
        title.className = 'gallery-item-title';
        title.textContent = item.title || '';
        
        // Description
        if (item.description) {
            const description = document.createElement('p');
            description.className = 'gallery-item-description';
            description.textContent = item.description;
            overlay.appendChild(description);
        }
        
        // Tags
        if (item.tags && item.tags.length > 0) {
            const tags = document.createElement('div');
            tags.className = 'gallery-item-tags';
            item.tags.forEach(tag => {
                const tagEl = document.createElement('span');
                tagEl.className = 'gallery-item-tag';
                tagEl.textContent = tag;
                tags.appendChild(tagEl);
            });
            overlay.appendChild(tags);
        }
        
        imgContainer.appendChild(img);
        imgContainer.appendChild(overlay);
        
        itemEl.appendChild(imgContainer);
        itemEl.onclick = () => this.options.onItemClick?.(item);
        
        return itemEl;
    }

    attachEventListeners() {
        // Intersection Observer for lazy loading
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    const img = entry.target;
                    img.src = img.dataset.src;
                    observer.unobserve(img);
                }
            });
        }, {
            rootMargin: '50px'
        });
        
        // Observe all gallery images
        this.container.querySelectorAll('.gallery-item img[data-src]').forEach(img => {
            observer.observe(img);
        });
    }

    updateItems(items, append = false) {
        if (append) {
            this.options.items = [...this.options.items, ...items];
        } else {
            this.options.items = items;
        }
        this.render();
    }

    setLoading(loading) {
        this.options.loading = loading;
        const loadMore = this.container.querySelector('.gallery-load-more');
        if (loadMore) {
            loadMore.textContent = loading ? 'Loading...' : 'Load More';
            loadMore.disabled = loading;
        }
    }

    setHasMore(hasMore) {
        this.options.hasMore = hasMore;
        this.render();
    }
}

// Export the Gallery class
window.Gallery = Gallery; 