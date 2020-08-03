import { Injectable, EventEmitter } from '@angular/core';

@Injectable({
    providedIn: 'root'
})
export class ThemeService {

    public themes = ['dashboard-light-theme', 'dashboard-dark-theme'];
    public themeIndex = 0;
    public themeChanged: EventEmitter<string> = new EventEmitter();

    constructor() { }

    getTheme(): string {
        return this.themes[this.themeIndex];
    }

    getThemes(): string[] {
        return this.themes;
    }

    changeTheme() {
        this.themeIndex = this.themeIndex + 1;
        if (this.themeIndex >= this.themes.length) { this.themeIndex = 0; }
        this.themeChanged.emit(this.getTheme());
    }
}

