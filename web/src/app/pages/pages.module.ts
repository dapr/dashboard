import { NgModule } from '@angular/core';
import { PagesComponent } from './pages.component';
import { DashboardModule } from './dashboard/dashboard.module';
import { PagesRoutingModule } from './pages-routing.module';
import { CommonModule } from '@angular/common';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { DetailModule } from './dashboard/detail/detail.module';
import { DaprComponentsModule } from './dapr-components/dapr-components.module';
import { MatListModule } from '@angular/material/list';
import { ConfigurationModule } from './configuration/configuration.module';
import { ControlPlaneModule } from './controlplane/controlplane.module';
import { DaprComponentDetailModule } from './dapr-components/dapr-component-detail/dapr-component-detail.module';
import { ConfigurationDetailModule } from './configuration/configuration-detail/configuration-detail.module';
import { OverlayModule } from '@angular/cdk/overlay';
import { MatSelectModule } from '@angular/material/select';
import { FormsModule } from '@angular/forms';
import { MatDialogModule } from '@angular/material/dialog';

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
    DaprComponentDetailModule,
    ConfigurationDetailModule,
    OverlayModule,
    MatSelectModule,
    FormsModule,
    MatDialogModule
  ],
  declarations: [
    PagesComponent,
  ],
})
export class PagesModule { }
