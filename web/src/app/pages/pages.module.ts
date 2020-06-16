import { NgModule } from '@angular/core';
import { NbMenuModule, NbCardModule } from '@nebular/theme';
import { ThemeModule } from '../@theme/theme.module';
import { PagesComponent } from './pages.component';
import { DashboardModule } from './dashboard/dashboard.module';
import { ComponentsModule } from './components/components.module';
import { PagesRoutingModule } from './pages-routing.module';
import { DetailModule } from './detail/detail.module';

@NgModule({
  imports: [
    PagesRoutingModule,
    ThemeModule,
    NbMenuModule,
    DashboardModule,
    ComponentsModule,
    DetailModule,
    NbCardModule,
  ],
  declarations: [
    PagesComponent,
  ],
})
export class PagesModule {}
