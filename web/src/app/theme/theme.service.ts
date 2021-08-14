import { Injectable, EventEmitter } from '@angular/core';

const STORAGE_THEME_KEY = 'preferred_dashboard_theme';

@Injectable({
    providedIn: 'root'
})
export class ThemeService {

    public themes = ['dashboard-light-theme', 'dashboard-dark-theme'];
    public themeIndex = 0;
    public themeChanged: EventEmitter<string> = new EventEmitter();

    constructor() { }

    getTheme(): string {
        const savedThemeItem = localStorage.getItem(STORAGE_THEME_KEY);
        const savedThemeIndex = parseInt(savedThemeItem, 10);

        if (!isNaN(savedThemeIndex) && savedThemeIndex < this.themes.length) {
            this.themeIndex = savedThemeIndex;
        }

        return this.themes[this.themeIndex];
    }

    getThemes(): string[] {
        return this.themes;
    }

    changeTheme() {
        this.themeIndex = this.themeIndex + 1;
        if (this.themeIndex >= this.themes.length) { this.themeIndex = 0; }
        this.themeChanged.emit(this.themes[this.themeIndex]);

        localStorage.setItem(STORAGE_THEME_KEY, `${this.themeIndex}`);
    }

    isDarkTheme() {
        return this.getTheme().includes('dark');
    }
}

