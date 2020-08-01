import { Injectable } from '@angular/core';

@Injectable({
    providedIn: 'root'
})
export class ThemeService {

    public themes = ['dashboard-light-theme', 'dashboard-dark-theme'];
    public themeIndex = 0;

    constructor() { }

    getTheme(): string {
        return this.themes[this.themeIndex];
    }

    changeTheme(): void {
        this.themeIndex = this.themeIndex + 1;
        if (this.themeIndex >= this.themes.length) this.themeIndex = 0;
    }
}

