/**
 * TimeBride Application JavaScript
 * Містить утилітарні функції та основні налаштування для додатку
 */

// Глобальний об'єкт з утилітами
window.TimeBride = window.TimeBride || {};

// Налаштування для всього додатку
TimeBride.config = {
    apiBaseUrl: '/api',
    dateFormat: 'DD.MM.YYYY',
    timeFormat: 'HH:mm',
    datetimeFormat: 'DD.MM.YYYY HH:mm',
    defaultLocale: 'uk',
    debounceTimeout: 300,
    modalAnimationDuration: 300
};

// Утиліти для роботи з DOM
TimeBride.dom = {
    /**
     * Показати модальне вікно з динамічним вмістом
     * @param {string} title - Заголовок модального вікна
     * @param {string|HTMLElement} content - HTML вміст або DOM-елемент
     * @param {Object} options - Додаткові налаштування
     */
    showModal: function(title, content, options = {}) {
        const modal = document.getElementById('globalModal');
        if (!modal) return;
        
        const modalTitle = modal.querySelector('.modal-title');
        const modalBody = modal.querySelector('.modal-body');
        
        // Встановлюємо заголовок
        if (modalTitle) {
            modalTitle.textContent = title;
        }
        
        // Очищаємо та встановлюємо новий вміст
        if (modalBody) {
            modalBody.innerHTML = '';
            
            if (typeof content === 'string') {
                modalBody.innerHTML = content;
            } else if (content instanceof HTMLElement) {
                modalBody.appendChild(content);
            }
        }
        
        // Встановлюємо розмір діалогу
        const dialog = modal.querySelector('.modal-dialog');
        if (dialog) {
            dialog.className = 'modal-dialog'; // Скидаємо класи
            if (options.size) {
                dialog.classList.add(`modal-${options.size}`);
            } else {
                dialog.classList.add('modal-lg'); // За замовчуванням
            }
            
            if (options.centered) {
                dialog.classList.add('modal-dialog-centered');
            }
            
            if (options.scrollable) {
                dialog.classList.add('modal-dialog-scrollable');
            }
        }
        
        // Показуємо модальне вікно
        const modalInstance = new bootstrap.Modal(modal);
        modalInstance.show();
        
        return modalInstance;
    },
    
    /**
     * Створює повідомлення сповіщення (toast)
     * @param {string} message - Текст повідомлення
     * @param {string} type - Тип повідомлення (success, danger, warning, info)
     * @param {Object} options - Додаткові налаштування
     */
    showToast: function(message, type = 'info', options = {}) {
        // Створюємо контейнер для сповіщень, якщо він не існує
        let toastContainer = document.querySelector('.toast-container');
        if (!toastContainer) {
            toastContainer = document.createElement('div');
            toastContainer.className = 'toast-container position-fixed top-0 end-0 p-3';
            document.body.appendChild(toastContainer);
        }
        
        // Час показу повідомлення
        const delay = options.delay || 5000;
        
        // ID для повідомлення
        const id = 'toast-' + Date.now();
        
        // Створюємо HTML для повідомлення
        const toastHtml = `
            <div id="${id}" class="toast" role="alert" aria-live="assertive" aria-atomic="true" data-bs-delay="${delay}">
                <div class="toast-header bg-${type} bg-opacity-10">
                    <i class="ti ${type === 'success' ? 'ti-circle-check' : 
                                    type === 'danger' ? 'ti-alert-circle' : 
                                    type === 'warning' ? 'ti-alert-triangle' : 
                                    'ti-info-circle'} text-${type} me-2"></i>
                    <strong class="me-auto">${options.title || 'Повідомлення'}</strong>
                    <small>${options.subtitle || ''}</small>
                    <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
                </div>
                <div class="toast-body">
                    ${message}
                </div>
            </div>
        `;
        
        // Додаємо повідомлення в контейнер
        toastContainer.insertAdjacentHTML('beforeend', toastHtml);
        
        // Отримуємо елемент і показуємо його
        const toastEl = document.getElementById(id);
        const toast = new bootstrap.Toast(toastEl);
        toast.show();
        
        // Видаляємо елемент з DOM після сховати
        toastEl.addEventListener('hidden.bs.toast', function() {
            toastEl.remove();
        });
        
        return toast;
    }
};

// Утиліти для роботи з даними
TimeBride.data = {
    /**
     * Форматує дату за вказаним форматом
     * @param {Date|string} date - Дата для форматування
     * @param {string} format - Формат дати (за замовчуванням з налаштувань)
     * @returns {string} Відформатована дата
     */
    formatDate: function(date, format = TimeBride.config.dateFormat) {
        if (!date) return '';
        
        // Якщо використовується dayjs або luxon, можна використати їх
        // Простий приклад форматування для DD.MM.YYYY
        const d = new Date(date);
        const day = d.getDate().toString().padStart(2, '0');
        const month = (d.getMonth() + 1).toString().padStart(2, '0');
        const year = d.getFullYear();
        
        return `${day}.${month}.${year}`;
    },
    
    /**
     * Форматує час за вказаним форматом
     * @param {Date|string} date - Дата для форматування часу
     * @param {string} format - Формат часу (за замовчуванням з налаштувань)
     * @returns {string} Відформатований час
     */
    formatTime: function(date, format = TimeBride.config.timeFormat) {
        if (!date) return '';
        
        // Простий приклад форматування для HH:mm
        const d = new Date(date);
        const hours = d.getHours().toString().padStart(2, '0');
        const minutes = d.getMinutes().toString().padStart(2, '0');
        
        return `${hours}:${minutes}`;
    },
    
    /**
     * Форматує дату і час за вказаним форматом
     * @param {Date|string} date - Дата для форматування
     * @param {string} format - Формат дати і часу (за замовчуванням з налаштувань)
     * @returns {string} Відформатована дата і час
     */
    formatDateTime: function(date, format = TimeBride.config.datetimeFormat) {
        if (!date) return '';
        
        // Простий приклад для DD.MM.YYYY HH:mm
        return TimeBride.data.formatDate(date) + ' ' + TimeBride.data.formatTime(date);
    }
};

// Утиліти для роботи з API
TimeBride.api = {
    /**
     * Виконує запит до API
     * @param {string} endpoint - Endpoint API
     * @param {Object} options - Налаштування запиту
     * @returns {Promise} Promise з результатом запиту
     */
    request: async function(endpoint, options = {}) {
        const url = `${TimeBride.config.apiBaseUrl}${endpoint}`;
        
        const defaultOptions = {
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            credentials: 'same-origin'
        };
        
        const fetchOptions = {...defaultOptions, ...options};
        
        if (options.body && typeof options.body === 'object') {
            fetchOptions.body = JSON.stringify(options.body);
        }
        
        try {
            const response = await fetch(url, fetchOptions);
            
            // Перевіряємо успішність запиту
            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.message || `API request failed with status ${response.status}`);
            }
            
            // Якщо у відповіді немає вмісту
            if (response.status === 204) {
                return null;
            }
            
            // Парсимо JSON відповідь
            return await response.json();
        } catch (error) {
            console.error('API request error:', error);
            
            // Показуємо повідомлення про помилку
            if (TimeBride.dom && TimeBride.dom.showToast) {
                TimeBride.dom.showToast(
                    error.message || 'Помилка запиту до сервера',
                    'danger',
                    { title: 'Помилка API' }
                );
            }
            
            throw error;
        }
    },
    
    // Методи-хелпери для різних типів запитів
    get: function(endpoint, params = {}) {
        const queryString = new URLSearchParams(params).toString();
        const url = queryString ? `${endpoint}?${queryString}` : endpoint;
        return this.request(url, { method: 'GET' });
    },
    
    post: function(endpoint, data = {}) {
        return this.request(endpoint, { 
            method: 'POST',
            body: data
        });
    },
    
    put: function(endpoint, data = {}) {
        return this.request(endpoint, { 
            method: 'PUT',
            body: data
        });
    },
    
    delete: function(endpoint) {
        return this.request(endpoint, { method: 'DELETE' });
    }
};

// Ініціалізація при завантаженні сторінки
document.addEventListener('DOMContentLoaded', function() {
    // Ініціалізація спливаючих підказок (tooltips)
    const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]');
    const tooltipList = [...tooltipTriggerList].map(tooltipTriggerEl => 
        new bootstrap.Tooltip(tooltipTriggerEl)
    );
    
    // Ініціалізація popover
    const popoverTriggerList = document.querySelectorAll('[data-bs-toggle="popover"]');
    const popoverList = [...popoverTriggerList].map(popoverTriggerEl => 
        new bootstrap.Popover(popoverTriggerEl)
    );
    
    // Інші ініціалізації...
    console.log('TimeBride application initialized');
}); 