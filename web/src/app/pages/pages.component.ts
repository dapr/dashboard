import { Component, ViewChild, OnInit, HostBinding } from '@angular/core';
import { MenuItem, MENU_ITEMS, COMPONENTS_MENU_ITEM, CONFIGURATIONS_MENU_ITEM, CONTROLPLANE_MENU_ITEM } from './pages-menu';
import { FeaturesService } from 'src/app/features/features.service';
import { MatSidenav } from '@angular/material/sidenav';
import { GlobalsService } from 'src/app/globals/globals.service';
import { ThemeService } from 'src/app/theme/theme.service';

@Component({
  selector: 'app-pages',
  styleUrls: ['pages.component.scss'],
  templateUrl: 'pages.component.html',
})
export class PagesComponent implements OnInit {

  @HostBinding('class') public componentCssClass: string;
  @ViewChild('drawer', { static: false })
  public drawer: MatSidenav;
  public menu: MenuItem[] = MENU_ITEMS;
  public isMenuOpen = false;
  public contentMargin = 60;
  public isLightMode = true;

  constructor(
    private featuresService: FeaturesService,
    public globals: GlobalsService,
    private themeService: ThemeService,
  ) { }

  ngOnInit(): void {
    this.getFeatures();
    this.componentCssClass = this.themeService.getTheme();
  }

  getFeatures(): void {
    this.featuresService.get().subscribe((data: string[]) => {
      for (const feature of data) {
        if (feature === COMPONENTS_MENU_ITEM.name) {
          this.menu.push(COMPONENTS_MENU_ITEM);
        }
        if (feature === CONFIGURATIONS_MENU_ITEM.name) {
          this.menu.push(CONFIGURATIONS_MENU_ITEM);
        }
        if (feature === CONTROLPLANE_MENU_ITEM.name) {
          this.menu.push(CONTROLPLANE_MENU_ITEM);
        }
      }
    });
  }

  onDrawerToggle(): void {
    this.isMenuOpen = !this.isMenuOpen;
    if (!this.isMenuOpen) {
      this.contentMargin = 60;
    }
    else {
      this.contentMargin = 240;
    }
  }

  onThemeChange(): void {
    this.themeService.changeTheme();
    this.componentCssClass = this.themeService.getTheme();
    this.isLightMode = !this.isLightMode;
  }
}
