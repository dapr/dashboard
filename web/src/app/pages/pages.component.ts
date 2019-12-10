import { Component } from '@angular/core';
import { MENU_ITEMS, COMPONENTS_MENU_ITEM } from './pages-menu';
import { FeaturesService } from '../features/features.service';

@Component({
  selector: 'ngx-pages',
  styleUrls: ['pages.component.scss'],
  template: `
    <ngx-one-column-layout>
      <nb-menu [items]="menu"></nb-menu>
      <router-outlet></router-outlet>
    </ngx-one-column-layout>
  `,
})
export class PagesComponent {
  menu = MENU_ITEMS;

  constructor(private features: FeaturesService) {
    this.getFeatures();
  }

  getFeatures() {
    this.features.get().subscribe((data: string[]) => {
      for (let feature of data) {
        if (feature == COMPONENTS_MENU_ITEM.name) {
          this.menu.push(COMPONENTS_MENU_ITEM)
        }
      }
    });
  }
}
