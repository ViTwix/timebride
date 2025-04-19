// Calendar Component
class Calendar {
    constructor(container, options = {}) {
        this.container = container;
        this.options = {
            view: options.view || 'month',
            date: options.date || new Date(),
            events: options.events || [],
            onEventClick: options.onEventClick || null,
            onDateClick: options.onDateClick || null,
            onViewChange: options.onViewChange || null
        };
        
        this.init();
    }

    init() {
        this.render();
        this.attachEventListeners();
    }

    render() {
        const calendar = document.createElement('div');
        calendar.className = 'calendar';
        
        // Header
        const header = this.createHeader();
        calendar.appendChild(header);
        
        // Calendar grid
        const grid = this.createGrid();
        calendar.appendChild(grid);
        
        this.container.innerHTML = '';
        this.container.appendChild(calendar);
    }

    createHeader() {
        const header = document.createElement('div');
        header.className = 'calendar-header';
        
        // Navigation buttons
        const prevBtn = document.createElement('button');
        prevBtn.className = 'calendar-nav-btn';
        prevBtn.innerHTML = '&lt;';
        prevBtn.onclick = () => this.navigate(-1);
        
        const nextBtn = document.createElement('button');
        nextBtn.className = 'calendar-nav-btn';
        nextBtn.innerHTML = '&gt;';
        nextBtn.onclick = () => this.navigate(1);
        
        // Current month/year
        const title = document.createElement('h2');
        title.className = 'calendar-title';
        title.textContent = this.getFormattedDate();
        
        // View selector
        const viewSelector = document.createElement('select');
        viewSelector.className = 'calendar-view-selector';
        ['month', 'week', 'day'].forEach(view => {
            const option = document.createElement('option');
            option.value = view;
            option.textContent = view.charAt(0).toUpperCase() + view.slice(1);
            option.selected = view === this.options.view;
            viewSelector.appendChild(option);
        });
        viewSelector.onchange = (e) => this.changeView(e.target.value);
        
        header.appendChild(prevBtn);
        header.appendChild(title);
        header.appendChild(nextBtn);
        header.appendChild(viewSelector);
        
        return header;
    }

    createGrid() {
        const grid = document.createElement('div');
        grid.className = `calendar-grid calendar-${this.options.view}`;
        
        if (this.options.view === 'month') {
            // Week days header
            const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
            const weekHeader = document.createElement('div');
            weekHeader.className = 'calendar-week-header';
            weekDays.forEach(day => {
                const dayEl = document.createElement('div');
                dayEl.className = 'calendar-week-day';
                dayEl.textContent = day;
                weekHeader.appendChild(dayEl);
            });
            grid.appendChild(weekHeader);
            
            // Days grid
            const daysGrid = document.createElement('div');
            daysGrid.className = 'calendar-days';
            
            const firstDay = new Date(this.options.date.getFullYear(), this.options.date.getMonth(), 1);
            const lastDay = new Date(this.options.date.getFullYear(), this.options.date.getMonth() + 1, 0);
            const startDay = firstDay.getDay();
            const totalDays = lastDay.getDate();
            
            // Previous month days
            for (let i = 0; i < startDay; i++) {
                const dayEl = document.createElement('div');
                dayEl.className = 'calendar-day calendar-day-prev';
                daysGrid.appendChild(dayEl);
            }
            
            // Current month days
            for (let i = 1; i <= totalDays; i++) {
                const dayEl = document.createElement('div');
                dayEl.className = 'calendar-day';
                dayEl.textContent = i;
                
                // Add events for this day
                const dayEvents = this.getEventsForDay(i);
                if (dayEvents.length > 0) {
                    const eventsContainer = document.createElement('div');
                    eventsContainer.className = 'calendar-day-events';
                    dayEvents.forEach(event => {
                        const eventEl = document.createElement('div');
                        eventEl.className = 'calendar-event';
                        eventEl.textContent = event.title;
                        eventEl.onclick = () => this.options.onEventClick?.(event);
                        eventsContainer.appendChild(eventEl);
                    });
                    dayEl.appendChild(eventsContainer);
                }
                
                dayEl.onclick = () => this.options.onDateClick?.(new Date(this.options.date.getFullYear(), this.options.date.getMonth(), i));
                daysGrid.appendChild(dayEl);
            }
            
            grid.appendChild(daysGrid);
        }
        
        return grid;
    }

    getEventsForDay(day) {
        return this.options.events.filter(event => {
            const eventDate = new Date(event.start_time);
            return eventDate.getDate() === day &&
                   eventDate.getMonth() === this.options.date.getMonth() &&
                   eventDate.getFullYear() === this.options.date.getFullYear();
        });
    }

    getFormattedDate() {
        const options = { month: 'long', year: 'numeric' };
        return this.options.date.toLocaleDateString('en-US', options);
    }

    navigate(direction) {
        if (this.options.view === 'month') {
            this.options.date.setMonth(this.options.date.getMonth() + direction);
        } else if (this.options.view === 'week') {
            this.options.date.setDate(this.options.date.getDate() + (direction * 7));
        } else {
            this.options.date.setDate(this.options.date.getDate() + direction);
        }
        this.render();
    }

    changeView(view) {
        this.options.view = view;
        this.options.onViewChange?.(view);
        this.render();
    }

    attachEventListeners() {
        // Add any additional event listeners here
    }

    updateEvents(events) {
        this.options.events = events;
        this.render();
    }
}

// Export the Calendar class
window.Calendar = Calendar; 