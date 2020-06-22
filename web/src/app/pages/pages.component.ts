import { Component, ViewChild } from '@angular/core';
import { MenuItem, MENU_ITEMS, COMPONENTS_MENU_ITEM } from './pages-menu';
import { FeaturesService } from '../features/features.service';
import { MatSidenav } from '@angular/material/sidenav';

@Component({
  selector: 'ngx-pages',
  styleUrls: ['pages.component.scss'],
  templateUrl: 'pages.component.html',
})
export class PagesComponent {
  menu: MenuItem[] = MENU_ITEMS;

  constructor(private features: FeaturesService) {
    this.getFeatures();
  }

  @ViewChild('drawer', { static: false }) 
  drawer: MatSidenav;

  getFeatures() {
    this.features.get().subscribe((data: string[]) => {
      for (const feature of data) {
        if (feature === COMPONENTS_MENU_ITEM.name) {
          this.menu.push(COMPONENTS_MENU_ITEM);
        }
      }
    });
  }
}
