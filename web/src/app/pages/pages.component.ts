import { OverlayContainer } from '@angular/cdk/overlay';
import { Component, ElementRef, HostBinding, Inject, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { FeaturesService } from 'src/app/features/features.service';
import { GlobalsService } from 'src/app/globals/globals.service';
import { ThemeService } from 'src/app/theme/theme.service';
import { VERSION } from '../../environments/version';
import { ScopesService } from '../scopes/scopes.service';
import { COMPONENTS_MENU_ITEM, CONFIGURATIONS_MENU_ITEM, CONTROLPLANE_MENU_ITEM, MenuItem, MENU_ITEMS } from './pages-menu';
import {DaprVersion} from "../types/types";

export interface DialogData {
  version: string;
  runtimeVersion: string;
}

@Component({
  selector: 'app-pages',
  styleUrls: ['pages.component.scss'],
  templateUrl: 'pages.component.html',
})
export class PagesComponent implements OnInit, OnDestroy {

  @HostBinding('class') public componentCssClass!: string;

  public menu: MenuItem[] = MENU_ITEMS;
  public isMenuOpen = false;
  public scopeValue = 'All';
  public scopes: string[] = [];
  public version: string | undefined;
  public runtimeVersion: string | undefined;
  public versionLoaded = false;
  private intervalHandler: any;

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
    this.applyTheme();

    this.intervalHandler = setInterval(() => {
      this.getScopes();
    }, 10000);
  }

  ngOnDestroy(): void {
    clearInterval(this.intervalHandler);
  }

  getVersion(): void {
    this.globals.getVersion().subscribe((version: DaprVersion) => {
      this.version = version.version;
      this.runtimeVersion = version.runtimeVersion;
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
  }

  onThemeChange(): void {
    this.themeService.changeTheme();
    this.applyTheme();
  }

  private applyTheme(): void {
    this.componentCssClass = this.themeService.getTheme();
    this.themeService.getThemes().forEach(theme => {
      this.overlayContainer.getContainerElement().classList.remove(theme);
    });
    this.overlayContainer.getContainerElement().classList.add(this.themeService.getTheme());
  }

  onScopeChange(): void {
    this.scopesService.changeScope(this.scopeValue);
    this.getScopes();
  }

  openDialog() {
    this.dialog.open(AboutDialogComponent, {
      data: {
        version: this.version,
        runtimeVersion: this.runtimeVersion
      } as DialogData
    });
  }

  isLightMode(): boolean {
    return this.componentCssClass === 'dashboard-light-theme';
  }
}

@Component({
  selector: 'app-about-dialog',
  templateUrl: 'dialog-template.html',
  styles: [`
    td button {
      visibility: hidden;
    }

    td:hover button {
      visibility: initial;
    }
  `]
})
export class AboutDialogComponent {
  @ViewChild('info', { static: true }) public info!: ElementRef;
  public version = VERSION;

  dashboardVersion = (this.version.semver as any)?.version || this.version.version;

  constructor(@Inject(MAT_DIALOG_DATA) public data: DialogData) { }

  copyInfo(data?: string) {
    const result = data || this.info.nativeElement?.innerText?.replace(/( )*content_copy( )*/g, '') || '';
    navigator.clipboard.writeText(result.replace(/\n\n/g, '\n'));
  }
}
