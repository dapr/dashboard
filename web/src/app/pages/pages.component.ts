import { Component, ViewChild, OnInit } from '@angular/core';
import { MenuItem, MENU_ITEMS, COMPONENTS_MENU_ITEM, CONFIGURATIONS_MENU_ITEM, CONTROLPLANE_MENU_ITEM } from './pages-menu';
import { FeaturesService } from '../features/features.service';
import { MatSidenav } from '@angular/material/sidenav';
import { GlobalsService } from '../globals/globals.service';

@Component({
  selector: 'ngx-pages',
  styleUrls: ['pages.component.scss'],
  templateUrl: 'pages.component.html',
})
export class PagesComponent implements OnInit {

  public menu: MenuItem[] = MENU_ITEMS;
  public isMenuOpen = false;
  public contentMargin = 60;

  constructor(
    private features: FeaturesService,
    public globals: GlobalsService,
  ) { }

  ngOnInit() {
    this.getFeatures();
    this.globals.getSupportedEnvironments().subscribe(data => {
      let supportedEnvironments = <Array<any>>data;
      if (supportedEnvironments.includes("kubernetes")) {
        this.globals.kubernetesEnabled = true;
      }
      else if (supportedEnvironments.includes("standalone")) {
        this.globals.standaloneEnabled = true;
      }
    });
  }

  @ViewChild('drawer', { static: false })
  drawer: MatSidenav;

  getFeatures() {
    this.features.get().subscribe((data: string[]) => {
      for (let feature of data) {
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

  onDrawerToggle() {
    this.isMenuOpen = !this.isMenuOpen;
    if (!this.isMenuOpen) {
      this.contentMargin = 60;
    }
    else {
      this.contentMargin = 240;
    }
  }
}
