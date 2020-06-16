import { NgModule } from '@angular/core';
import { PagesComponent } from './pages.component';
import { DashboardModule } from './dashboard/dashboard.module';
import { PagesRoutingModule } from './pages-routing.module';
import { CommonModule } from '@angular/common';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';

@NgModule({
  imports: [
    CommonModule,
    PagesRoutingModule,
    DashboardModule,
    MatSidenavModule,
    MatToolbarModule,
    MatIconModule,
  ],
  declarations: [
    PagesComponent,
  ],
})
export class PagesModule {}
