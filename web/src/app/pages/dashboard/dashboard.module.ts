import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NbCardModule, NbLayoutModule, NbIconModule, NbDialogModule } from '@nebular/theme';
import { ThemeModule } from '../../@theme/theme.module';
import { DashboardComponent } from './dashboard.component';

@NgModule({
  imports: [
    NbCardModule,
    ThemeModule,
    NbLayoutModule,
    NbIconModule,
    RouterModule,
    NbDialogModule.forRoot()
  ],
  declarations: [
    DashboardComponent,
  ]
})
export class DashboardModule { }
