import { NgModule } from '@angular/core';
import { NbMenuModule, NbCardModule } from '@nebular/theme';

import { ThemeModule } from '../@theme/theme.module';
import { PagesComponent } from './pages.component';
import { DashboardModule } from './dashboard/dashboard.module';
import { ComponentsModule } from './components/components.module';
import { LogsModule } from './logs/logs.module';
import { PagesRoutingModule } from './pages-routing.module';
import { ExtraOptions, RouterModule, Routes, Router } from '@angular/router';

@NgModule({
  imports: [
    PagesRoutingModule,
    ThemeModule,
    NbMenuModule,
    DashboardModule,
    ComponentsModule,
    LogsModule,
    NbCardModule,
    RouterModule
  ],
  declarations: [
    PagesComponent,
  ],
})
export class PagesModule {
}
