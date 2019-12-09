import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NbCardModule, NbLayoutModule, NbIconModule, NbDialogModule } from '@nebular/theme';

import { ThemeModule } from '../../@theme/theme.module';
import { DashboardComponent } from './dashboard.component';
import { DataTableModule } from 'ng-angular8-datatable';

@NgModule({
  imports: [
    NbCardModule,
    ThemeModule,
    DataTableModule,
    NbLayoutModule,
    NbIconModule,
    NbDialogModule,
    RouterModule
  ],
  declarations: [
    DashboardComponent,
  ]
})
export class DashboardModule { }
