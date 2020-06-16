import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { MatTableModule } from '@angular/material/table';

@NgModule({
  imports: [
    RouterModule,
    MatTableModule,
  ],
  declarations: [
    DashboardComponent,
  ],
})
export class DashboardModule { }
