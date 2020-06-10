import { NgModule } from '@angular/core';
import { NbMenuModule, NbCardModule, NbThemeModule, NbLayoutModule, NbSidebarModule } from '@nebular/theme';

import { PagesComponent } from './pages.component';
import { DashboardModule } from './dashboard/dashboard.module';
import { ComponentsModule } from './components/components.module';
import { LogsModule } from './logs/logs.module';
import { PagesRoutingModule } from './pages-routing.module';
import { ThemeModule } from '../@theme/theme.module';

@NgModule({
  imports: [
    PagesRoutingModule,
    NbMenuModule,
    DashboardModule,
    ComponentsModule,
    LogsModule,
    NbCardModule,
    NbLayoutModule,
    NbSidebarModule,
    NbThemeModule,
    ThemeModule
  ],
  declarations: [
    PagesComponent,
  ],
})
export class PagesModule {
}
