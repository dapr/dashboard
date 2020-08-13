import { Component, ViewChild, OnInit, HostBinding, OnDestroy, Inject } from '@angular/core';
import { MenuItem, MENU_ITEMS, COMPONENTS_MENU_ITEM, CONFIGURATIONS_MENU_ITEM, CONTROLPLANE_MENU_ITEM } from './pages-menu';
import { FeaturesService } from 'src/app/features/features.service';
import { MatSidenav } from '@angular/material/sidenav';
import { GlobalsService } from 'src/app/globals/globals.service';
import { ThemeService } from 'src/app/theme/theme.service';
import { OverlayContainer } from '@angular/cdk/overlay';
import { Router } from '@angular/router';
import { ScopesService } from '../scopes/scopes.service';
import { MatDialog, MAT_DIALOG_DATA } from '@angular/material/dialog';

export interface DialogData {
  version: string;
}

@Component({
  selector: 'app-pages',
  styleUrls: ['pages.component.scss'],
  templateUrl: 'pages.component.html',
})
export class PagesComponent implements OnInit, OnDestroy {

  @HostBinding('class') public componentCssClass: string;
  @ViewChild('drawer', { static: false })
  public drawer: MatSidenav;
  public menu: MenuItem[] = MENU_ITEMS;
  public isMenuOpen = false;
  public contentMargin = 60;
  public isLightMode = true;
  public imgPath: string;
  public themeSelectorEnabled: boolean;
  public scopeValue = 'All';
  public scopes: string[];
  public version: string;
  public versionLoaded: boolean;
  private intervalHandler;

  constructor(
    private featuresService: FeaturesService,
    public globals: GlobalsService,
    private themeService: ThemeService,
    private overlayContainer: OverlayContainer,
    public router: Router,
    private scopesService: ScopesService,
    public dialog: MatDialog,
  ) { }

  ngOnInit(): void {
    this.getVersion();
    this.getFeatures();
    this.getScopes();
    this.componentCssClass = this.themeService.getTheme();
    this.imgPath = '../../assets/images/logo.svg';

    this.intervalHandler = setInterval(() => {
      this.getScopes();
    }, 10000);
  }

  ngOnDestroy(): void {
    clearInterval(this.intervalHandler);
  }

  getVersion(): void {
    this.globals.getVersion().subscribe(version => {
      this.version = version;
      this.versionLoaded = true;
    });
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

  getScopes(): void {
    this.scopesService.getScopes().subscribe((data: string[]) => {
      this.scopes = data;
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
    this.themeService.getThemes().forEach(theme => {
      this.overlayContainer.getContainerElement().classList.remove(theme);
    });
    this.overlayContainer.getContainerElement().classList.add(this.themeService.getTheme());
    if (this.isLightMode) {
      this.imgPath = '../../assets/images/logo.svg';
    } else {
      this.imgPath = '../../assets/images/logo-white.svg';
    }
  }

  onScopeChange(): void {
    this.scopesService.changeScope(this.scopeValue);
    this.getScopes();
  }

  openDialog() {
    this.dialog.open(AboutDialogComponent, {
      data: {
        version: this.version
      }
    });
  }
}

@Component({
  selector: 'app-about-dialog',
  templateUrl: 'dialog-template.html',
})
export class AboutDialogComponent {
  constructor(@Inject(MAT_DIALOG_DATA) public data: DialogData) {}
}
