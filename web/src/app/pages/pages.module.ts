import { NgModule } from '@angular/core';
import { PagesComponent } from './pages.component';
import { DashboardModule } from './dashboard/dashboard.module';
import { PagesRoutingModule } from './pages-routing.module';
import { CommonModule } from '@angular/common';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { DetailModule } from './detail/detail.module';
import { DaprComponentsModule } from './dapr-components/dapr-components.module';
import { MatListModule } from '@angular/material/list';
import { ConfigurationModule } from './configuration/configuration.module';
import { ControlPlaneModule } from './controlplane/controlplane.module';

@NgModule({
  imports: [
    CommonModule,
    PagesRoutingModule,
    DashboardModule,
    DaprComponentsModule,
    DetailModule,
    MatSidenavModule,
    MatToolbarModule,
    MatIconModule,
    MatListModule,
    ConfigurationModule,
    ControlPlaneModule,
  ],
  declarations: [
    PagesComponent,
  ],
})
export class PagesModule { }
